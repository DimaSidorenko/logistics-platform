package dto

import "errors"

var (
	ErrUserNotFound = errors.New("cart: user not found")
	ErrItemNotFound = errors.New("cart: item not found")
)
