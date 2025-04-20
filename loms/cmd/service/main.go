package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/kafka"
	"route256/loms/internal/infra/kafka/producer"
	"route256/loms/internal/logger"
	"route256/loms/internal/middlewares"
	lomsService "route256/loms/internal/service/loms"
	lomsUsecase "route256/loms/internal/usecases/loms"
	"route256/loms/internal/usecases/loms/stocks_repository"
	lomsUsecaseStorage "route256/loms/internal/usecases/loms/storage"
	desc "route256/loms/pkg/protobuf/rpc/server"
)

//go:embed stock-data.json
var stockData embed.FS

func initStockData() ([]lomsUsecaseStorage.Stock, error) {
	// Чтение заэмбеденного файла
	data, err := stockData.ReadFile("stock-data.json")
	if err != nil {
		fmt.Println("Error reading embedded file:", err)
		return nil, err
	}

	// Парсинг JSON
	var stocks []lomsUsecaseStorage.Stock
	err = json.Unmarshal(data, &stocks)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return stocks, nil
}

func main() {
	ctx := context.Background()

	stocks, err := initStockData()
	if err != nil {
		panic(err)
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	jaegerResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("loms"),
		),
	)
	if err != nil {
		panic(err)
	}

	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://localhost:4318"))
	if err != nil {
		panic(err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(jaegerResource),
	)

	defer func() {
		_ = traceProvider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(traceProvider)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Service.GrpcPort))
	if err != nil {
		panic(err)
	}

	dbDsnMaster := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		"postgresql",          // протокол
		cfg.DbMaster.User,     // имя пользователя
		cfg.DbMaster.Password, // пароль
		cfg.DbMaster.Host,     // хост
		cfg.DbMaster.Port,     // порт
		cfg.DbMaster.DbName,   // имя базы данных
		"disable",             // параметр sslmode
	)

	masterConfig, err := pgxpool.ParseConfig(dbDsnMaster)
	if err != nil {
		log.Fatalf("unable to parse config: %v\n", err)
	}

	masterPool, err := pgxpool.NewWithConfig(ctx, masterConfig)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v\n", err)
	}

	// Init loms.order-events kafka producer.
	orderEventsProducer, err := producer.NewSyncProducer(
		kafka.Config{Brokers: []string{cfg.Kafka.Brokers}},
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		//producer.WithProducerPartitioner(sarama.NewManualPartitioner),
		//producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		//producer.WithProducerPartitioner(sarama.NewRandomPartitioner),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
	)
	if err != nil {
		log.Fatalf("unable to create loms order events producer: %v\n", err)
	}

	usecase := lomsUsecase.NewUsecase(
		stocks_repository.NewRepositoryDB(masterPool),
		lomsUsecaseStorage.NewStocksStorage(stocks),
		orderEventsProducer,
		cfg.Kafka.OrderTopic,
	)
	service := lomsService.NewService(usecase)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // конструкция для трейсинга
		grpc.ChainUnaryInterceptor(
			middlewares.Logger,
			middlewares.Validate,
			middlewares.Metrics,
			//middlewares.Tracing,
		),
	)
	reflection.Register(grpcServer)
	desc.RegisterLomsServer(grpcServer, service)

	go func() {
		logger.Warnw(ctx, "Serving grpcServer on %s\n", list.Addr().String())
		if err = grpcServer.Serve(list); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", cfg.Service.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	// Init router for http server.
	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(
			func(header string) (string, bool) {
				switch strings.ToLower(header) {
				case "x-auth":
					return header, true
				default:
					return header, true
				}
			},
		),
	)
	if err = desc.RegisterLomsHandler(ctx, gwmux, conn); err != nil {
		panic(err)
	}

	metricsHandler := promhttp.Handler()
	if err = gwmux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		metricsHandler.ServeHTTP(w, r)
	}); err != nil {
		log.Fatalf("Error serving /metrics: %v\n", err)
	}

	//nolint:gosec
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Service.HttpPort),
		Handler: gwmux,
	}

	logger.Warnw(ctx, "Serving gRPC-Gateway on %s\n", gwServer.Addr)
	if err = gwServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
