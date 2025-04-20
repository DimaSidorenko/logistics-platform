package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"route256/cart/internal/usecases/cart/dto"
)

func TestGetItem(t *testing.T) {
	// Arrange.
	storage := NewConcurrentMap()
	userID := dto.UserID(10)
	skuID1 := dto.SkuID(101)
	skuID2 := dto.SkuID(102)

	_ = storage.AddItem(userID, skuID1, 5)

	t.Run("existing item", func(t *testing.T) {
		// Action.
		quantity, found := storage.GetItem(userID, skuID1)
		assert.True(t, found)
		assert.Equal(t, uint32(5), quantity)
	})

	t.Run("non-existing item", func(t *testing.T) {
		// Action
		quantity, found := storage.GetItem(userID, skuID2)
		assert.False(t, found)
		assert.Equal(t, uint32(0), quantity)
	})

	t.Run("non-existing user", func(t *testing.T) {
		// Action.
		quantity, found := storage.GetItem(dto.UserID(999), skuID1)
		assert.False(t, found)
		assert.Equal(t, uint32(0), quantity)
	})
}

func TestGetItems(t *testing.T) {
	storage := NewConcurrentMap()
	userID := dto.UserID(10)
	skuID1 := dto.SkuID(101)
	skuID2 := dto.SkuID(102)

	err := storage.AddItem(userID, skuID1, 5)
	require.NoError(t, err)
	err = storage.AddItem(userID, skuID2, 10)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("existing user", func(t *testing.T) {
		items, err := storage.GetItems(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, items, 2)
		assert.ElementsMatch(t, items, []dto.Item{{Sku: skuID1, Count: 5}, {Sku: skuID2, Count: 10}})
	})

	t.Run("non-existing user", func(t *testing.T) {
		items, err := storage.GetItems(ctx, dto.UserID(999))
		require.NoError(t, err)
		assert.Len(t, items, 0)
	})
}

func TestAddItem(t *testing.T) {
	storage := NewConcurrentMap()
	userID := dto.UserID(1)
	skuID := dto.SkuID(101)

	err := storage.AddItem(userID, skuID, 3)
	require.NoError(t, err)
	quantity1, found1 := storage.GetItem(userID, skuID)

	_ = storage.AddItem(userID, skuID, 2)
	require.NoError(t, err)
	quantity2, _ := storage.GetItem(userID, skuID)

	require.True(t, found1)
	assert.Equal(t, uint32(3), quantity1)
	assert.Equal(t, uint32(5), quantity2)
}

func TestDeleteItem(t *testing.T) {
	storage := NewConcurrentMap()
	userID := dto.UserID(1)
	skuID1 := dto.SkuID(101)
	skuID2 := dto.SkuID(102)

	_ = storage.AddItem(userID, skuID1, 5)
	//nolint:errcheck
	storage.AddItem(userID, skuID2, 3)

	t.Run("delete existing item", func(t *testing.T) {
		err := storage.DeleteItem(userID, skuID1)
		require.NoError(t, err)
		quantity, found := storage.GetItem(userID, skuID1)
		assert.False(t, found)
		assert.Equal(t, uint32(0), quantity)

		quantity, found = storage.GetItem(userID, skuID2)
		assert.True(t, found)
		assert.Equal(t, uint32(3), quantity)
	})

	t.Run("delete last item should remove user", func(t *testing.T) {
		err := storage.DeleteItem(userID, skuID2)
		require.NoError(t, err)
		_, found := storage.storage[userID]
		assert.False(t, found)
	})

	t.Run("delete non-existing item", func(t *testing.T) {
		err := storage.DeleteItem(userID, skuID1)
		require.NoError(t, err)
		_, found := storage.GetItem(userID, skuID1)
		assert.False(t, found)
	})
}

func TestDeleteUser(t *testing.T) {
	storage := NewConcurrentMap()
	userID := dto.UserID(1)
	skuID := dto.SkuID(101)

	err := storage.AddItem(userID, skuID, 5)
	require.NoError(t, err)
	err = storage.DeleteUser(userID)
	require.NoError(t, err)
	quantity, found := storage.GetItem(userID, skuID)
	assert.False(t, found)
	assert.Equal(t, uint32(0), quantity)
}
