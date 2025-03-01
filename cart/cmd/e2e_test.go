package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"route256/cart/internal/app/handlers/add_item"
	"route256/cart/internal/app/handlers/delete_item"
	"route256/cart/internal/app/handlers/delete_user"
	"route256/cart/internal/app/handlers/get_items"
	"route256/cart/internal/infra/http/middlewares"
	"route256/cart/internal/infra/http/round_trippers"
	"route256/cart/internal/services/product"
	"route256/cart/internal/usecases/cart"
	"route256/cart/internal/usecases/cart/repository"
)

const (
	token               = "testToken"
	productServiceImage = "gitlab-registry.ozon.dev/go/classroom-16/students/base/products:latest"
)

type E2ESuite struct {
	suite.Suite
	container testcontainers.Container
	server    *http.Server
}

func (p *E2ESuite) SetupTest() {
	p.T().Log("Starting product service container")
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        productServiceImage,
		ExposedPorts: []string{"8082/tcp"},
		WaitingFor:   wait.ForHTTP("/docs").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	p.Require().NoError(err, "Ошибка при создании контейнера")
	p.container = container

	endpoint, err := p.container.Endpoint(ctx, "")
	p.Require().NoError(err)
	p.server = p.startCartServer("http://" + endpoint)
}

func (p *E2ESuite) TearDownTest() {
	p.T().Log("Stopping product service container")
	ctx := context.Background()
	if p.container != nil {
		err := p.container.Terminate(ctx)
		p.Require().NoError(err, "Ошибка при остановке контейнера")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if p.server != nil {
		if err := p.server.Shutdown(ctx); err != nil {
			p.T().Fatalf("Server Shutdown Failed: %+v", err)
		}
	}
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

func (p *E2ESuite) addItem(userID int64, skuID int64) int {
	data := add_item.AddItemRequest{
		Count: 10,
	}

	jsonData, err := json.Marshal(data)
	p.Require().NoError(err, "marshal JSON")
	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/user/%v/cart/%v", userID, skuID), "application/json", bytes.NewBuffer(jsonData))
	p.Require().NoError(err, "add item http request")

	return resp.StatusCode
}

func (p *E2ESuite) getUserCart(userID int64) (statusCode int, result *get_items.GetItemsResponse) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/user/%v/cart", userID))
	p.Require().NoError(err, "get items http request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	p.Require().NoError(err, "read response body")
	statusCode = resp.StatusCode

	if statusCode != http.StatusOK {
		return statusCode, nil
	}

	err = json.Unmarshal(body, &result)
	p.Require().NoError(err, "unmarshal get items http request")
	return
}

func (p *E2ESuite) deleteUserItem(userID, skuID int64) int {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/user/%v/cart/%v", userID, skuID), nil)
	p.Require().NoError(err, "create delete request")

	client := &http.Client{}
	resp, err := client.Do(req)
	p.Require().NoError(err, "send delete request")
	defer resp.Body.Close()

	return resp.StatusCode
}

func (p *E2ESuite) TestDeleteItem() {
	userID := int64(1)
	itemID := int64(1076963)

	code := p.addItem(userID, itemID)
	p.Assert().Equal(http.StatusOK, code)

	code, userCart := p.getUserCart(userID)
	p.Equal(http.StatusOK, code)
	p.Len(userCart.Items, 1)

	code = p.deleteUserItem(userID, itemID)
	p.Assert().Equal(http.StatusNoContent, code)

	code, _ = p.getUserCart(userID)
	p.Assert().Equal(http.StatusNotFound, code)
}

func (p *E2ESuite) TestListItem() {
	userID := int64(1)
	item1ID := int64(1076963)
	item2ID := int64(1625903)

	code, _ := p.getUserCart(userID)
	p.Assert().Equal(http.StatusNotFound, code)

	code = p.addItem(userID, item1ID)
	p.Assert().Equal(http.StatusOK, code)

	code, userCart := p.getUserCart(userID)
	p.Equal(http.StatusOK, code)
	p.Len(userCart.Items, 1)

	code = p.addItem(userID, item2ID)
	p.Assert().Equal(http.StatusOK, code)

	code, userCart = p.getUserCart(userID)
	p.Equal(http.StatusOK, code)
	p.Len(userCart.Items, 2)

	p.Equal(userCart.TotalPrice, uint32(48020))

	sort.Slice(userCart.Items, func(i, j int) bool {
		return userCart.Items[i].Sku < userCart.Items[j].Sku
	})

	p.Equal(userCart.Items[0].Sku, item1ID)
	p.Equal(userCart.Items[0].Count, uint32(10))
	p.Equal(userCart.Items[1].Sku, item2ID)
	p.Equal(userCart.Items[1].Count, uint32(10))
}

func (p *E2ESuite) startCartServer(productServiceURL string) *http.Server {
	productClient := product.NewProductClient(
		&http.Client{Transport: round_trippers.NewRetryRoundTripper(http.DefaultTransport, 3)},
		productServiceURL,
		token,
	)

	cartHandler := cart.NewHandler(productClient, repository.NewConcurrentMap())
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

	go func() {
		log.Println("Starting cartServer on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Printf("cartServer finished: %v\n", err)
		}
	}()

	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", "localhost:8080", 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return server
		}
		time.Sleep(500 * time.Millisecond)
	}

	p.T().Fatalf("cartServer failed to start after")
	return server
}
