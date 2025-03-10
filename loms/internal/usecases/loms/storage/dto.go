package storage

import (
	"route256/loms/internal/usecases/loms/dto"
	"time"
)

type Order struct {
	ID        int64
	UserID    int64
	Status    dto.OrderStatus
	Items     []dto.Item
	CreatedAt time.Time
	UpdatedAt time.Time
}
