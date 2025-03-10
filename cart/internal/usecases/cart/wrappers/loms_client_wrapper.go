package wrappers

import (
	"context"

	"google.golang.org/grpc"

	cartDto "route256/cart/internal/usecases/cart/dto"
	desc "route256/cart/pkg/protobuf/rpc/clients"
)

type LomsClient interface {
	OrderCreate(ctx context.Context, in *desc.OrderCreateRequest, opts ...grpc.CallOption) (*desc.OrderCreateResponse, error)
	StocksInfo(ctx context.Context, in *desc.StocksInfoRequest, opts ...grpc.CallOption) (*desc.StocksInfoResponse, error)
}

type LomsClientWrapper struct {
	client LomsClient
}

func NewLomsClientWrapper(client LomsClient) *LomsClientWrapper {
	return &LomsClientWrapper{client: client}
}

func (w *LomsClientWrapper) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	req := &desc.StocksInfoRequest{
		Sku: sku,
	}

	resp, err := w.client.StocksInfo(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (w *LomsClientWrapper) OrderCreate(ctx context.Context, userID cartDto.UserID, items []cartDto.Item) (orderID int64, err error) {
	var orderItems []*desc.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, &desc.OrderItem{
			SKU:   int64(item.Sku),
			Count: item.Count,
		})
	}

	req := &desc.OrderCreateRequest{
		User:  int64(userID),
		Items: orderItems,
	}

	resp, err := w.client.OrderCreate(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.OrderID, nil
}
