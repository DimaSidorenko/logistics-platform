package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"route256/loms/internal/usecases/loms/dto"
)

func TestStockRepository(t *testing.T) {
	ctx := context.Background()
	storage := NewStocksStorage(nil)

	t.Run("ReserveStocks success", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 123, TotalCount: 100, Reserved: 0})

		items := []dto.Item{
			{SKU: 123, Count: 10},
		}

		err := storage.ReserveStocks(ctx, items)
		assert.NoError(t, err)

		stock, _ := storage.GetAvailableStock(ctx, 123)
		assert.Equal(t, uint32(90), stock)
	})

	t.Run("ReserveStocks failure", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 456, TotalCount: 100, Reserved: 0})

		items := []dto.Item{
			{SKU: 456, Count: 150},
		}

		err := storage.ReserveStocks(ctx, items)
		assert.Error(t, err)
		assert.Equal(t, ErrNotEnoughStock, err)
	})

	t.Run("CancelReservation success", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 789, TotalCount: 100, Reserved: 20})

		err := storage.CancelReservation(ctx, 789, 10)
		assert.NoError(t, err)

		stock, _ := storage.GetAvailableStock(ctx, 789)
		assert.Equal(t, uint32(90), stock)
	})

	t.Run("CancelReservation failure", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 999, TotalCount: 100, Reserved: 10})

		err := storage.CancelReservation(ctx, 999, 20)
		assert.Error(t, err)
		assert.Equal(t, ErrReservationNotFound, err)
	})

	t.Run("RemoveReservation success", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 111, TotalCount: 100, Reserved: 10})

		err := storage.RemoveReservation(ctx, 111, 5)
		assert.NoError(t, err)

		stock, _ := storage.GetAvailableStock(ctx, 111)
		assert.Equal(t, uint32(90), stock)
	})

	t.Run("RemoveReservation failure", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 222, TotalCount: 100, Reserved: 10})

		err := storage.RemoveReservation(ctx, 222, 20)
		assert.Error(t, err)
		assert.Equal(t, ErrReservationNotFound, err)
	})

	t.Run("GetAvailableStock success", func(t *testing.T) {
		storage.AddStock(Stock{SKU: 333, TotalCount: 100, Reserved: 30})

		available, err := storage.GetAvailableStock(ctx, 333)
		assert.NoError(t, err)
		assert.Equal(t, uint32(70), available)
	})

	t.Run("GetAvailableStock failure", func(t *testing.T) {
		_, err := storage.GetAvailableStock(ctx, 444)
		assert.Error(t, err)
	})
}
