package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
	config, err := config2.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	productClient := product.NewProductClient(
		&http.Client{Transport: round_trippers.NewRetryRoundTripper(http.DefaultTransport, 3)},
		fmt.Sprintf("http://%s:%d", config.ProductService.Host, config.ProductService.Port),
		config.ProductService.Token,
	)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", config.LomsService.Host, config.LomsService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	lomsClient := desc.NewLomsClient(conn)

	cartHandler := cart.NewHandler(productClient, wrappers.NewLomsClientWrapper(lomsClient), repository.NewConcurrentMap())

	mux := http.NewServeMux()

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", add_item.NewHandler(cartHandler))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", delete_item.NewHandler(cartHandler))
	mux.Handle("DELETE /user/{user_id}/cart", delete_user.NewHandler(cartHandler))
	mux.Handle("GET /user/{user_id}/cart", get_items.NewHandler(cartHandler))
	mux.Handle("POST /checkout/{user_id}", checkout_order.NewHandler(cartHandler))

	server := &http.Server{
		Addr:        ":8080",
		Handler:     middlewares.NewLoggingMiddleware(mux),
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
