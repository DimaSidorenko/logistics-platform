package main

import (
	"fmt"
	"log"
	"net/http"
	"route256/cart/internal/app/handlers/add_item"
	"route256/cart/internal/app/handlers/delete_item"
	"route256/cart/internal/app/handlers/delete_user"
	"route256/cart/internal/app/handlers/get_items"
	config2 "route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/middlewares"
	"route256/cart/internal/infra/http/round_trippers"
	"route256/cart/internal/services"
	"route256/cart/internal/usecases/cart"
	"time"
)

func main() {
	config, err := config2.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	productClient := services.NewProductClient(
		&http.Client{Transport: round_trippers.NewRetryRoundTripper(http.DefaultTransport, 3)},
		config.ProductService.Host,
		config.ProductService.Port,
		config.ProductService.Token,
	)

	cartHandler := cart.NewHandler(productClient)

	mux := http.NewServeMux()

	mux.Handle("POST /user/{user_id}/cart/{sku_id}", add_item.NewHandler(cartHandler))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id}", delete_item.NewHandler(cartHandler))
	mux.Handle("DELETE /user/{user_id}/cart", delete_user.NewHandler(cartHandler))
	mux.Handle("GET /user/{user_id}/cart", get_items.NewHandler(cartHandler))

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
