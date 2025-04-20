package repository

import (
	"context"
	"route256/cart/internal/tracing"
	"sync"

	cartDto "route256/cart/internal/usecases/cart/dto"
)

type ConcurrentMap struct {
	mu      sync.RWMutex
	storage map[cartDto.UserID]map[cartDto.SkuID]uint32
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		storage: make(map[cartDto.UserID]map[cartDto.SkuID]uint32),
	}
}

func (c *ConcurrentMap) GetItem(userID cartDto.UserID, skuID cartDto.SkuID) (uint32, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userCart, exists := c.storage[userID]
	if !exists {
		return 0, false
	}
	quantity, found := userCart[skuID]
	return quantity, found
}

func (c *ConcurrentMap) GetItems(ctx context.Context, userID cartDto.UserID) ([]cartDto.Item, error) {
	_, span := tracing.StartFromContext(ctx, "concurrentMap.GetItems")
	defer span.End()

	c.mu.RLock()
	defer c.mu.RUnlock()

	userCart, exists := c.storage[userID]
	if !exists {
		return nil, nil
	}

	items := make([]cartDto.Item, 0, len(userCart))
	for skuID, quantity := range userCart {
		items = append(items, cartDto.Item{Sku: skuID, Count: quantity})
	}
	return items, nil
}

func (c *ConcurrentMap) AddItem(userID cartDto.UserID, skuID cartDto.SkuID, quantity uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.storage[userID]; !exists {
		c.storage[userID] = make(map[cartDto.SkuID]uint32)
	}
	c.storage[userID][skuID] += quantity
	return nil
}

func (c *ConcurrentMap) DeleteItem(userID cartDto.UserID, skuID cartDto.SkuID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	userCart, exists := c.storage[userID]
	if !exists {
		return nil
	}
	delete(userCart, skuID)
	if len(userCart) == 0 {
		delete(c.storage, userID)
	}
	return nil
}

func (c *ConcurrentMap) DeleteUser(userID cartDto.UserID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.storage, userID)
	return nil
}
