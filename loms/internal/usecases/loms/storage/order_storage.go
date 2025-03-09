package storage

import (
	"context"
	"sync"
	"time"

	"route256/loms/internal/usecases/loms/dto"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[int64]*Order
	nextID int64
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[int64]*Order),
		nextID: 1,
	}
}

func (s *OrderStorage) CreateOrder(_ context.Context, userID int64, items []dto.Item) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &Order{
		ID:        s.nextID,
		UserID:    userID,
		Status:    dto.OrderStatusNew,
		Items:     items,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.orders[order.ID] = order
	s.nextID++

	return order.ID, nil
}

func (s *OrderStorage) UpdateOrderStatus(_ context.Context, orderID int64, status dto.OrderStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
	if !exists {
		return dto.ErrOrderNotFound
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

func (s *OrderStorage) GetOrderByID(_ context.Context, orderID int64) (*dto.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[orderID]
	if !exists {
		return nil, dto.ErrOrderNotFound
	}

	return &dto.Order{
		User:   order.UserID,
		Status: order.Status,
		Items:  order.Items,
	}, nil
}
