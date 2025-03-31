package loms

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"

	"route256/loms/internal/usecases/loms/dto"
	"route256/loms/internal/usecases/loms/storage"
)

var _ OrderRepository = (*storage.OrderStorage)(nil)

type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int64, items []dto.Item) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status dto.OrderStatus) error
	GetOrderByID(ctx context.Context, orderID int64) (*dto.Order, error)
}

var _ StockRepository = (*storage.StocksStorage)(nil)

type StockRepository interface {
	ReserveStocks(ctx context.Context, items []dto.Item) error
	RemoveReservation(ctx context.Context, sku int64, count uint32) error
	CancelReservation(ctx context.Context, sku int64, count uint32) error
	GetAvailableStock(ctx context.Context, sku int64) (uint32, error)
}

type orderEventProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type Usecase struct {
	orderRepo OrderRepository
	stockRepo StockRepository

	orderEventProducer orderEventProducer
	topic              string
}

func NewUsecase(orderRepo OrderRepository, stockRepo StockRepository, orderEventProducer orderEventProducer, topic string) *Usecase {
	return &Usecase{
		orderRepo:          orderRepo,
		stockRepo:          stockRepo,
		orderEventProducer: orderEventProducer,
		topic:              topic,
	}
}
func (u *Usecase) produceOrderEvent(orderID string, message string) {
	msg := sarama.ProducerMessage{
		Topic:     u.topic,
		Key:       sarama.StringEncoder(orderID),
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
	}
	_, _, _ = u.orderEventProducer.SendMessage(&msg)
}

func (u *Usecase) OrderCreate(ctx context.Context, userID int64, items []dto.Item) (int64, error) {
	orderID, err := u.orderRepo.CreateOrder(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	u.produceOrderEvent(strconv.FormatInt(orderID, 10), "create new order")

	if err = u.stockRepo.ReserveStocks(ctx, items); err != nil {
		log.Printf("failed to reserve stocks: %v", err)
		_ = u.orderRepo.UpdateOrderStatus(ctx, orderID, dto.OrderStatusFailed)

		u.produceOrderEvent(strconv.FormatInt(orderID, 10), "reserve stocks for order failed")

		return 0, dto.ErrReserveFailed
	}

	if err = u.orderRepo.UpdateOrderStatus(ctx, orderID, dto.OrderStatusAwaitingPayment); err != nil {
		return 0, fmt.Errorf("update order status %v", err)
	}

	u.produceOrderEvent(strconv.FormatInt(orderID, 10), "awaiting payment")

	return orderID, nil
}

func (u *Usecase) OrderInfo(ctx context.Context, orderID int64) (*dto.Order, error) {
	order, err := u.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, dto.ErrOrderNotFound
	}

	return order, nil
}

func (u *Usecase) OrderPay(ctx context.Context, orderID int64) error {
	order, err := u.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("get order by id: %w", err)
	}

	if order.Status == dto.OrderStatusPayed {
		return nil
	}

	if order.Status == dto.OrderStatusCancelled {
		return dto.ErrOrderCancelled
	}

	if order.Status != dto.OrderStatusAwaitingPayment {
		return dto.ErrOrderNotAwaitingPayment
	}

	for _, item := range order.Items {
		if err = u.stockRepo.RemoveReservation(ctx, item.SKU, item.Count); err != nil {
			return err
		}
	}

	if err = u.orderRepo.UpdateOrderStatus(ctx, orderID, dto.OrderStatusPayed); err != nil {
		return err
	}

	u.produceOrderEvent(strconv.FormatInt(orderID, 10), "order payed")

	return nil
}

func (u *Usecase) OrderCancel(ctx context.Context, orderID int64) error {
	order, err := u.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return dto.ErrOrderNotFound
	}

	if order.Status == dto.OrderStatusCancelled {
		return nil
	}

	if order.Status == dto.OrderStatusPayed {
		return dto.ErrCannotCancelOrder
	}

	for _, item := range order.Items {
		err = u.stockRepo.CancelReservation(ctx, item.SKU, item.Count)
		if err != nil {
			return err
		}
	}

	err = u.orderRepo.UpdateOrderStatus(ctx, orderID, dto.OrderStatusCancelled)
	if err != nil {
		return err
	}

	u.produceOrderEvent(strconv.FormatInt(orderID, 10), "order canceled")

	return nil
}

func (u *Usecase) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	count, err := u.stockRepo.GetAvailableStock(ctx, sku)
	if err != nil {
		return 0, err
	}

	return count, nil
}
