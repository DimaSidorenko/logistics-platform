package cart

import (
	"errors"
	dto2 "route256/cart/internal/usecases/cart/dto"
)

var storage = map[dto2.UserID]map[dto2.SkuID]uint32{}

type productClient interface {
	GetItem(skuID dto2.SkuID) (dto2.Product, error)
}

type Handler struct {
	productClient productClient
}

func NewHandler(productClient productClient) *Handler {
	return &Handler{
		productClient: productClient,
	}
}

func (c *Handler) AddItem(userID dto2.UserID, skuID dto2.SkuID, quantity uint32) error {
	if err := c.validateItem(skuID); err != nil {
		return err
	}

	if userCart, ok := storage[userID]; ok {
		userCart[skuID] = userCart[skuID] + quantity
	} else {
		storage[userID] = map[dto2.SkuID]uint32{
			skuID: quantity,
		}
	}

	return nil
}

func (c *Handler) DeleteItem(userID dto2.UserID, skuID dto2.SkuID) error {
	if userCart, ok := storage[userID]; ok {
		delete(userCart, skuID)
	}

	return nil
}

func (c *Handler) DeleteUser(userID dto2.UserID) error {
	delete(storage, userID)

	return nil
}

func (c *Handler) GetItems(userID dto2.UserID) (dto2.GetItemsResponse, error) {
	if _, ok := storage[userID]; !ok {
		return dto2.GetItemsResponse{}, dto2.ErrUserNotFound
	}

	result := make([]dto2.Item, 0, len(storage))
	var totalPrice uint32
	for skuID, quantity := range storage[userID] {
		item, err := c.productClient.GetItem(skuID)
		if err != nil {
			if errors.Is(err, dto2.ErrItemNotFound) {
				continue
			}

			return dto2.GetItemsResponse{}, err
		}

		price := uint32(item.Price)
		result = append(result, dto2.Item{
			Sku:   int64(skuID),
			Name:  item.Name,
			Count: quantity,
			Price: price,
		})

		totalPrice += price * quantity
	}

	return dto2.GetItemsResponse{
		Items:      result,
		TotalPrice: totalPrice,
	}, nil
}

func (c *Handler) validateItem(skuID dto2.SkuID) error {
	_, err := c.productClient.GetItem(skuID)
	return err
}
