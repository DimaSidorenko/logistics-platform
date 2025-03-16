package dto

import "errors"

var (
	ErrUserNotFound          = errors.New("cart: user not found")
	ErrFailedToReserveStocks = errors.New("cart: failed to reserve stocks")
)
