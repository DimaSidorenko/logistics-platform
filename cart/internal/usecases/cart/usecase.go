//go:generate minimock -i=storage
//go:generate minimock -i=productClient

package cart

import (
	"fmt"
	"route256/cart/internal/models"
	cartDto "route256/cart/internal/usecases/cart/dto"
)

type storage interface {
	GetItem(userID cartDto.UserID, skuID cartDto.SkuID) (quantity uint32, found bool)
	AddItem(userID cartDto.UserID, skuID cartDto.SkuID, quantity uint32) error
	DeleteItem(userID cartDto.UserID, skuID cartDto.SkuID) error
	DeleteUser(userID cartDto.UserID) error
	GetItems(userID cartDto.UserID) ([]cartDto.Item, error)
}

type productClient interface {
	GetItem(skuID int64) (models.Product, error)
}

type Handler struct {
	productClient productClient
	storage       storage
}

func NewHandler(productClient productClient, storage storage) *Handler {
	return &Handler{
		productClient: productClient,
		storage:       storage,
	}
}

func (c *Handler) AddItem(userID cartDto.UserID, skuID cartDto.SkuID, quantity uint32) error {
	if err := c.validateItem(skuID); err != nil {
		return fmt.Errorf("validate item: %v", err)
	}

	if err := c.storage.AddItem(userID, skuID, quantity); err != nil {
		return fmt.Errorf("add item: %v", err)
	}

	return nil
}

func (c *Handler) DeleteItem(userID cartDto.UserID, skuID cartDto.SkuID) error {
	if err := c.storage.DeleteItem(userID, skuID); err != nil {
		return fmt.Errorf("delete item: %v", err)
	}

	return nil
}

func (c *Handler) DeleteUser(userID cartDto.UserID) error {
	if err := c.storage.DeleteUser(userID); err != nil {
		return fmt.Errorf("delete item: %v", err)
	}

	return nil
}

func (c *Handler) GetItems(userID cartDto.UserID) (cartDto.GetItemsResponse, error) {
	items, err := c.storage.GetItems(userID)
	if err != nil {
		return cartDto.GetItemsResponse{}, fmt.Errorf("get items: %v", err)
	}

	if len(items) == 0 {
		return cartDto.GetItemsResponse{}, cartDto.ErrUserNotFound
	}

	result := make([]cartDto.Item, 0, len(items))
	var totalPrice uint32
	for _, repoItem := range items {
		item, err := c.productClient.GetItem(int64(repoItem.Sku))
		if err != nil {
			return cartDto.GetItemsResponse{}, fmt.Errorf("get item in product service: %v", err)
		}

		price := uint32(item.Price)
		result = append(result, cartDto.Item{
			Sku:   repoItem.Sku,
			Count: repoItem.Count,
			Name:  item.Name,
			Price: price,
		})
		totalPrice += price * repoItem.Count
	}

	return cartDto.GetItemsResponse{
		Items:      result,
		TotalPrice: totalPrice,
	}, nil
}

func (c *Handler) validateItem(skuID cartDto.SkuID) error {
	_, err := c.productClient.GetItem(int64(skuID))
	return err
}
