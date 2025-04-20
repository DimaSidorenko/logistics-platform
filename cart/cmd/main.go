package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/cart/internal/api/handlers/add_item"
	"route256/cart/internal/api/handlers/checkout_order"
	"route256/cart/internal/api/handlers/delete_item"
	"route256/cart/internal/api/handlers/delete_user"
	"route256/cart/internal/api/handlers/get_items"
	config2 "route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/middlewares"
	"route256/cart/internal/infra/http/round_trippers"
	"route256/cart/internal/services/product"
	"route256/cart/internal/usecases/cart"
	"route256/cart/internal/usecases/cart/repository"
	"route256/cart/internal/usecases/cart/wrappers"
	desc "route256/cart/pkg/protobuf/rpc/clients"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := config2.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	productClient := product.NewProductClient(
		&http.Client{Transport: round_trippers.NewRetryRoundTripper(http.DefaultTransport, 3)},
		fmt.Sprintf("http://%s:%d", config.ProductService.Host, config.ProductService.Port),
		config.ProductService.Token,
		rate.NewLimiter(10, 10),
		//ratelimiter.New(10),
	)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", config.LomsService.Host, config.LomsService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // Добавляет заголовки трейсинга
	)
	if err != nil {
		panic(err)
	}
	lomsClient := desc.NewLomsClient(conn)

	cartHandler := cart.NewHandler(productClient, wrappers.NewLomsClientWrapper(lomsClient), repository.NewConcurrentMap())

	r := mux.NewRouter()

	r.Handle("/user/{user_id}/cart/{sku_id}", add_item.NewHandler(cartHandler)).Methods("POST")
	r.Handle("/user/{user_id}/cart/{sku_id}", delete_item.NewHandler(cartHandler)).Methods("DELETE")
	r.Handle("/user/{user_id}/cart", delete_user.NewHandler(cartHandler)).Methods("DELETE")
	r.Handle("/user/{user_id}/cart", get_items.NewHandler(cartHandler)).Methods("GET")
	r.Handle("/checkout/{user_id}", checkout_order.NewHandler(cartHandler)).Methods("POST")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	r.Use(middlewares.NewMetricsMiddleware, middlewares.NewLoggingMiddleware)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     r,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		fmt.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error starting server: %v\n", err)
		}
	}()

	// Ожидаем отмены контекста (получения сигнала)
	<-ctx.Done()

	// Создаём контекст с таймаутом для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Пытаемся корректно завершить работу сервера
	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	//<-time.NewTimer(time.Second * 1).C

	log.Println("Server successfully shutdown")
}
