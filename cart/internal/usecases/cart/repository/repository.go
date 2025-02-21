package repository

import (
	"route256/cart/internal/usecases/cart/dto"
	"sync"
)

type ConcurrentMap struct {
	mu      sync.RWMutex
	storage map[dto.UserID]map[dto.SkuID]uint32
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		storage: make(map[dto.UserID]map[dto.SkuID]uint32),
	}
}

func (m *ConcurrentMap) Set(userID dto.UserID, skuID dto.SkuID, quantity uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.storage[userID]; !ok {
		m.storage[userID] = make(map[dto.SkuID]uint32)
	}
	m.storage[userID][skuID] += quantity
}

func (m *ConcurrentMap) DeleteItem(userID dto.UserID, skuID dto.SkuID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if userCart, ok := m.storage[userID]; ok {
		delete(userCart, skuID)
	}
}

func (m *ConcurrentMap) DeleteUser(userID dto.UserID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.storage, userID)
}

func (m *ConcurrentMap) Get(userID dto.UserID) (map[dto.SkuID]uint32, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userCart, ok := m.storage[userID]
	if !ok {
		return nil, false
	}

	cartCopy := make(map[dto.SkuID]uint32, len(userCart))
	for k, v := range userCart {
		cartCopy[k] = v
	}

	return cartCopy, true
}
