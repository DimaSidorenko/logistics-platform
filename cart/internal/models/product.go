package models

import "errors"

type Product struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Sku   int64  `json:"sku"`
}

var ErrItemNotFound = errors.New("cart: item not found")
