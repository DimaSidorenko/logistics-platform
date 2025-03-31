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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/kafka"
	"route256/loms/internal/infra/kafka/producer"
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
		grpc.ChainUnaryInterceptor(
			middlewares.Logger,
			middlewares.Validate,
		),
	)
	reflection.Register(grpcServer)
	desc.RegisterLomsServer(grpcServer, service)

	go func() {
		log.Printf("Serving grpcServer on %s\n", list.Addr().String())
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

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(
			func(header string) (string, bool) {
				switch strings.ToLower(header) {
				case "x-auth":
					return header, true
				default:
					return header, false
				}
			},
		),
	)
	if err = desc.RegisterLomsHandler(ctx, gwmux, conn); err != nil {
		panic(err)
	}

	//nolint:gosec
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Service.HttpPort),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	if err = gwServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
