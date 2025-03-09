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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/loms/internal/infra/config"
	"route256/loms/internal/middlewares"
	lomsService "route256/loms/internal/service/loms"
	lomsUsecase "route256/loms/internal/usecases/loms"
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

	usecase := lomsUsecase.NewUsecase(
		lomsUsecaseStorage.NewOrderStorage(),
		lomsUsecaseStorage.NewStocksStorage(stocks),
	)
	service := lomsService.NewService(usecase)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
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
