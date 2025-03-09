package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"route256/loms/internal/usecases/loms/dto"
)

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	orderStorage := NewOrderStorage()

	t.Run("CreateOrder success", func(t *testing.T) {
		items := []dto.Item{
			{SKU: 123, Count: 2},
			{SKU: 456, Count: 1},
		}

		orderID, err := orderStorage.CreateOrder(ctx, 1, items)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), orderID)
	})

	t.Run("UpdateOrderStatus success", func(t *testing.T) {
		items := []dto.Item{
			{SKU: 123, Count: 2},
			{SKU: 456, Count: 1},
		}

		orderID, _ := orderStorage.CreateOrder(ctx, 1, items)

		err := orderStorage.UpdateOrderStatus(ctx, orderID, "awaiting payment")
		assert.NoError(t, err)

		order, _ := orderStorage.GetOrderByID(ctx, orderID)
		assert.Equal(t, dto.OrderStatusAwaitingPayment, order.Status)
	})

	t.Run("UpdateOrderStatus failure", func(t *testing.T) {
		err := orderStorage.UpdateOrderStatus(ctx, 999, "awaiting payment")
		assert.Error(t, err)
		assert.Equal(t, dto.ErrOrderNotFound, err)
	})

	t.Run("GetOrderByID success", func(t *testing.T) {
		items := []dto.Item{
			{SKU: 123, Count: 2},
			{SKU: 456, Count: 1},
		}

		orderID, _ := orderStorage.CreateOrder(ctx, 1, items)

		order, err := orderStorage.GetOrderByID(ctx, orderID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), order.User)
		assert.Equal(t, dto.OrderStatusNew, order.Status)
		assert.Equal(t, items, order.Items)
	})

	t.Run("GetOrderByID failure", func(t *testing.T) {
		_, err := orderStorage.GetOrderByID(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, dto.ErrOrderNotFound, err)
	})
}
