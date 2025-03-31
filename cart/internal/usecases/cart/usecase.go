//go:generate minimock -i=storage
//go:generate minimock -i=productClient
//go:generate minimock -i=lomsClient

package cart

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/cart/internal/models"
	cartDto "route256/cart/internal/usecases/cart/dto"
	"route256/cart/internal/usecases/cart/wrappers"
)

type storage interface {
	GetItem(userID cartDto.UserID, skuID cartDto.SkuID) (quantity uint32, found bool)
	AddItem(userID cartDto.UserID, skuID cartDto.SkuID, quantity uint32) error
	DeleteItem(userID cartDto.UserID, skuID cartDto.SkuID) error
	DeleteUser(userID cartDto.UserID) error
	GetItems(userID cartDto.UserID) ([]cartDto.Item, error)
}

type productClient interface {
	GetItem(ctx context.Context, skuID int64) (models.Product, error)
}

var _ lomsClient = (*wrappers.LomsClientWrapper)(nil)

type lomsClient interface {
	StocksInfo(ctx context.Context, sku int64) (uint32, error)
	OrderCreate(context.Context, cartDto.UserID, []cartDto.Item) (orderID int64, err error)
}

type Handler struct {
	productClient productClient
	lomsClient    lomsClient
	storage       storage
}

func NewHandler(productClient productClient, lomsClient lomsClient, storage storage) *Handler {
	return &Handler{
		productClient: productClient,
		lomsClient:    lomsClient,
		storage:       storage,
	}
}

func (c *Handler) Checkout(userID cartDto.UserID) (orderID int64, err error) {
	resp, err := c.GetItems(userID)
	if err != nil {
		return 0, err
	}

	orderID, err = c.lomsClient.OrderCreate(context.TODO(), userID, resp.Items)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.FailedPrecondition {
				return 0, cartDto.ErrFailedToReserveStocks
			}
		}

		return 0, err
	}

	if err = c.DeleteUser(userID); err != nil {
		return 0, fmt.Errorf("delete user: %v", err)
	}

	return orderID, nil
}

// AddItem добавляет предметы в корзину.
func (c *Handler) AddItem(userID cartDto.UserID, skuID cartDto.SkuID, quantity uint32) error {
	if err := c.validateItem(skuID); err != nil {
		return fmt.Errorf("validate item: %v", err)
	}

	// (dosidorenko): по ходу курса менторы отказались от этого метода.
	//if err := c.validateProductExists(skuID, quantity); err != nil {
	//	return fmt.Errorf("validate product exists: %v", err)
	//}

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
	var mutex sync.Mutex

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	for _, repoItem := range items {
		eg.Go(func() error {
			item, err := c.productClient.GetItem(ctx, int64(repoItem.Sku))
			if err != nil {
				cancel()
				return fmt.Errorf("get item in product service: %v", err)
			}

			price := uint32(item.Price)

			func() {
				mutex.Lock()
				defer mutex.Unlock()

				result = append(result, cartDto.Item{
					Sku:   repoItem.Sku,
					Count: repoItem.Count,
					Name:  item.Name,
					Price: price,
				})
				totalPrice += price * repoItem.Count
			}()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return cartDto.GetItemsResponse{}, fmt.Errorf("get items: %v", err)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Sku < result[j].Sku
	})

	return cartDto.GetItemsResponse{
		Items:      result,
		TotalPrice: totalPrice,
	}, nil
}

func (c *Handler) validateItem(skuID cartDto.SkuID) error {
	_, err := c.productClient.GetItem(context.TODO(), int64(skuID))
	return err
}

// (dosidorenko): В конце курса можно будет удалить.
//func (c *Handler) validateProductExists(skuID cartDto.SkuID, neededCount uint32) error {
//	actualCount, err := c.lomsClient.StocksInfo(context.TODO(), int64(skuID))
//	if err != nil {
//		return fmt.Errorf("stocks info %v", err)
//	}
//
//	if actualCount < neededCount {
//		return fmt.Errorf("not enough stocks")
//	}
//
//	return nil
//}
