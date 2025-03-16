package stocks_repository

import (
	"context"
	"errors"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"route256/loms/internal/usecases/loms/dto"
)

type RepositoryDB struct {
	read  *pgxpool.Pool
	write *pgxpool.Pool
}

func NewRepositoryDB(pool *pgxpool.Pool) *RepositoryDB {
	return &RepositoryDB{read: pool, write: pool}
}

func (r RepositoryDB) CreateOrder(ctx context.Context, userID int64, items []dto.Item) (orderID int64, err error) {
	err = pgx.BeginTxFunc(ctx, r.write, pgx.TxOptions{}, func(tx pgx.Tx) (err error) {
		repo := New(tx)

		orderID, err = repo.CreateOrder(ctx, &CreateOrderParams{
			UserID: userID,
			Status: "new",
		})
		if err != nil {
			return err
		}

		for _, i := range items {
			err = repo.InsertOrderItem(ctx, &InsertOrderItemParams{
				OrderID:    orderID,
				SkuID:      i.SKU,
				ItemsCount: int64(i.Count),
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}

func (r RepositoryDB) UpdateOrderStatus(ctx context.Context, orderID int64, status dto.OrderStatus) (err error) {
	repo := New(r.write)
	err = repo.UpdateOrderStatus(ctx, &UpdateOrderStatusParams{
		ID:     orderID,
		Status: string(status),
	})

	return
}

func (r RepositoryDB) GetOrderByID(ctx context.Context, orderID int64) (order *dto.Order, err error) {
	err = pgx.BeginTxFunc(ctx, r.write, pgx.TxOptions{}, func(tx pgx.Tx) (err error) {
		repo := New(tx)

		orderInfo, err := repo.GetOrderInfo(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return dto.ErrOrderNotFound
			}
			return err
		}

		items, err := repo.GetOrderItems(ctx, orderID)
		if err != nil {
			return err
		}

		order = &dto.Order{
			OrderID: orderID,
			Status:  dto.OrderStatus(orderInfo.Status),
			User:    orderInfo.UserID,
			Items:   make([]dto.Item, 0, len(items)),
		}

		for _, item := range items {
			order.Items = append(order.Items, dto.Item{
				SKU: item.SkuID,
				// TODO(dosidorenko): make Count int64.
				//nolint:gosec
				Count: uint32(item.ItemsCount),
			})
		}

		return nil
	})

	if nil == err {
		sort.Slice(order.Items, func(i, j int) bool {
			return order.Items[i].SKU < order.Items[j].SKU
		})
	}

	return
}
