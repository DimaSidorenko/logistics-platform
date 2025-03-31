package stocks_repository

//
//import (
//	"context"
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/jackc/pgx/v5/pgxpool"
//	"github.com/jackc/pgx/v5/stdlib"
//	"github.com/pressly/goose/v3"
//	"github.com/stretchr/testify/suite"
//	"github.com/testcontainers/testcontainers-go"
//	"github.com/testcontainers/testcontainers-go/wait"
//
//	"route256/loms/internal/usecases/loms/dto"
//)
//
//type RepositoryDBTestSuite struct {
//	suite.Suite
//	pgContainer testcontainers.Container
//	pgEndpoint  string
//	dbName      string
//	dbUser      string
//	dbPassword  string
//	dbPool      *pgxpool.Pool
//}
//
//func (s *RepositoryDBTestSuite) SetupSuite() {
//	s.dbUser = "user"
//	s.dbPassword = "password"
//	s.dbName = "test_db"
//	const (
//		pgImage        = "gitlab-registry.ozon.dev/go/classroom-16/students/base/postgres:16"
//		migrationsPath = "../../../../migrations"
//	)
//
//	// Шаг 1: Поднимаем контейнер с PostgreSQL
//	ctx := context.Background()
//	postgresContainerRequest := testcontainers.ContainerRequest{
//		Image:        pgImage,
//		ExposedPorts: []string{"5432/tcp"},
//		Env: map[string]string{
//			"POSTGRES_USER":     s.dbUser,
//			"POSTGRES_PASSWORD": s.dbPassword,
//			"POSTGRES_DB":       s.dbName,
//		},
//		WaitingFor: wait.ForLog("database system is ready to accept connections").
//			WithOccurrence(1).
//			WithStartupTimeout(20 * time.Second),
//	}
//
//	var err error
//	s.pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
//		ContainerRequest: postgresContainerRequest,
//		Started:          true,
//	})
//	s.Require().NoError(err)
//
//	// Получаем endpoint контейнера
//	s.pgEndpoint, err = s.pgContainer.Endpoint(ctx, "")
//	s.Require().NoError(err)
//
//	// Шаг 2: Подключаемся к базе данных с использованием pgxpool.
//	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", s.dbUser, s.dbPassword, s.pgEndpoint, s.dbName)
//	fmt.Println(dsn)
//	s.dbPool, err = pgxpool.New(ctx, dsn)
//	s.Require().NoError(err, "Failed to connect to the database using pgxpool")
//
//	// Преобразуем pgx.Conn в стандартный интерфейс database/sql.
//	err = goose.Up(stdlib.OpenDBFromPool(s.dbPool), migrationsPath)
//	s.Require().NoError(err, "Failed to apply migrations")
//
//}
//
//func (s *RepositoryDBTestSuite) TearDownSuite() {
//	// Останавливаем контейнер после завершения тестов
//	ctx := context.Background()
//	s.Require().NoError(s.pgContainer.Terminate(ctx))
//}
//
//func (s *RepositoryDBTestSuite) TestUpdateOrderStatus() {
//	// Arrange.
//	ctx := context.Background()
//	initialStatus := dto.OrderStatusNew
//	updatedStatus := dto.OrderStatusPayed
//
//	querier := New(s.dbPool)
//
//	orderID, err := querier.CreateOrder(ctx, &CreateOrderParams{
//		Status: string(initialStatus),
//		UserID: 100,
//	})
//	s.Require().NoError(err, "Failed to insert test order")
//
//	// Action.
//	repo := NewRepositoryDB(s.dbPool)
//	err = repo.UpdateOrderStatus(ctx, orderID, updatedStatus)
//	s.Require().NoError(err, "Failed to update order status")
//
//	// Assert.
//	order, err := querier.GetOrderInfo(ctx, orderID)
//	s.Require().NoError(err, "Failed to query updated order status")
//	s.Equal(string(updatedStatus), order.Status, "get not updated Order status from db")
//}
//
//func (s *RepositoryDBTestSuite) TestCreateOrder() {
//	// Arrange.
//	ctx := context.Background()
//	userID := int64(100)
//	items := []dto.Item{
//		{SKU: 1, Count: 2},
//		{SKU: 2, Count: 3},
//	}
//
//	// Action.
//	repo := NewRepositoryDB(s.dbPool)
//	orderID, err := repo.CreateOrder(ctx, userID, items)
//	s.Require().NoError(err, "Failed to create order")
//
//	// Assert.
//	// Проверяем, что заказ создан
//	querier := New(s.dbPool)
//	order, err := querier.GetOrderInfo(ctx, orderID)
//	s.Require().NoError(err, "Failed to get order info")
//	s.Equal(userID, order.UserID, "UserID mismatch")
//	s.Equal(string(dto.OrderStatusNew), order.Status, "Order status mismatch")
//
//	// Проверяем, что товары добавлены в заказ
//	orderItems, err := querier.GetOrderItems(ctx, orderID)
//	s.Require().NoError(err, "Failed to get order items")
//	s.Equal(len(items), len(orderItems), "Number of items mismatch")
//
//	for i, item := range items {
//		s.Equal(item.SKU, orderItems[i].SkuID, "SKU mismatch")
//		s.Equal(int64(item.Count), orderItems[i].ItemsCount, "Item count mismatch")
//	}
//}
//
//func (s *RepositoryDBTestSuite) TestGetOrderByID() {
//	// Arrange.
//	ctx := context.Background()
//	userID := int64(100)
//	items := []dto.Item{
//		{SKU: 1, Count: 2},
//		{SKU: 2, Count: 3},
//	}
//
//	// Создаем заказ и добавляем товары
//	repo := NewRepositoryDB(s.dbPool)
//	orderID, err := repo.CreateOrder(ctx, userID, items)
//	s.Require().NoError(err, "Failed to create order")
//
//	// Action.
//	order, err := repo.GetOrderByID(ctx, orderID)
//	s.Require().NoError(err, "Failed to get order by ID")
//
//	// Assert.
//	// Проверяем, что заказ содержит правильные данные
//	s.Equal(orderID, order.OrderID, "OrderID mismatch")
//	s.Equal(dto.OrderStatusNew, order.Status, "Order status mismatch")
//	s.Equal(userID, order.User, "UserID mismatch")
//
//	// Проверяем, что товары в заказе соответствуют ожидаемым
//	s.Equal(len(items), len(order.Items), "Number of items mismatch")
//	for i, item := range items {
//		s.Equal(item.SKU, order.Items[i].SKU, "SKU mismatch")
//		s.Equal(item.Count, order.Items[i].Count, "Item count mismatch")
//	}
//}
//
//func TestRepositoryDBTestSuite(t *testing.T) {
//	suite.Run(t, new(RepositoryDBTestSuite))
//}
