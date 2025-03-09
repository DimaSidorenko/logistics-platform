package storage

import (
	"context"
	"errors"
	"sync"

	"route256/loms/internal/usecases/loms/dto"
)

var (
	ErrNotEnoughStock      = errors.New("not enough stock available")
	ErrReservationNotFound = errors.New("reservation not found")
)

type Stock struct {
	SKU        int64 `json:"sku" yaml:"sku"`
	TotalCount int   `json:"total_count" yaml:"total_count"` // Всего единиц товара на складе
	Reserved   int   `json:"reserved" yaml:"reserved"`       // Количество единиц в резерве
}

type StocksStorage struct {
	mu     sync.RWMutex
	stocks map[int64]*Stock
}

func NewStocksStorage(stocks []Stock) *StocksStorage {
	initMap := make(map[int64]*Stock)
	for _, stock := range stocks {
		initMap[stock.SKU] = &Stock{
			SKU:        stock.SKU,
			TotalCount: stock.TotalCount,
			Reserved:   stock.Reserved,
		}
	}

	return &StocksStorage{
		stocks: initMap,
	}
}

func (s *StocksStorage) AddStock(stock Stock) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stocks[stock.SKU] = &stock
}

// ReserveStocks резервирует items, если есть достаточное количество на складе.
func (s *StocksStorage) ReserveStocks(_ context.Context, items []dto.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, можно ли зарезервировать все товары
	for _, item := range items {
		stock, exists := s.stocks[item.SKU]
		if !exists {
			return ErrNotEnoughStock
		}
		if stock.TotalCount-stock.Reserved < int(item.Count) {
			return ErrNotEnoughStock
		}
	}

	// Резервируем товары
	for _, item := range items {
		stock := s.stocks[item.SKU]
		stock.Reserved += int(item.Count)
	}

	return nil
}

func (s *StocksStorage) RemoveReservation(_ context.Context, sku int64, count uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stock, exists := s.stocks[sku]
	if !exists || stock.Reserved < int(count) || stock.TotalCount < int(count) {
		return ErrReservationNotFound
	}

	stock.Reserved -= int(count)
	stock.TotalCount -= int(count)
	return nil
}

func (s *StocksStorage) CancelReservation(_ context.Context, sku int64, count uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stock, exists := s.stocks[sku]
	if !exists || stock.Reserved < int(count) {
		return ErrReservationNotFound
	}

	stock.Reserved -= int(count)
	return nil
}

func (s *StocksStorage) GetAvailableStock(_ context.Context, sku int64) (uint32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stock, exists := s.stocks[sku]
	if !exists {
		return 0, errors.New("sku not found")
	}

	available := stock.TotalCount - stock.Reserved
	if available < 0 {
		return 0, ErrNotEnoughStock
	}

	//nolint:gosec
	return uint32(available), nil
}
