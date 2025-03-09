package dto

import "errors"

var (
	ErrReserveFailed           = errors.New("failed to reserve stocks")
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderCancelled          = errors.New("order is cancelled")
	ErrOrderNotAwaitingPayment = errors.New("order is not awaiting payment")
	ErrCannotCancelOrder       = errors.New("cannot cancel order")
)
