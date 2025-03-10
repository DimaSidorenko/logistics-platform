package loms

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/loms/internal/usecases/loms"
	"route256/loms/internal/usecases/loms/dto"
	desc "route256/loms/pkg/protobuf/rpc/server"
)

var _ desc.LomsServer = (*Service)(nil)

var _ lomsUsecase = (*loms.Usecase)(nil)

type lomsUsecase interface {
	OrderCreate(ctx context.Context, userID int64, items []dto.Item) (orderID int64, err error)
	OrderInfo(ctx context.Context, orderID int64) (*dto.Order, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
	StocksInfo(ctx context.Context, sku int64) (uint32, error)
}

type Service struct {
	desc.UnimplementedLomsServer
	usecase lomsUsecase
}

func NewService(usecase lomsUsecase) *Service {
	return &Service{
		usecase: usecase,
	}
}

func (s *Service) OrderInfo(ctx context.Context, req *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	order, err := s.usecase.OrderInfo(ctx, req.OrderID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	resp := desc.OrderInfoResponse{
		Status: string(order.Status),
		User:   order.User,
		Items:  convertToOrderItems(order.Items),
	}

	return &resp, nil
}

func (s *Service) OrderCreate(ctx context.Context, req *desc.OrderCreateRequest) (*desc.OrderCreateResponse, error) {
	items := convertToDtoItems(req.Items)
	orderID, err := s.usecase.OrderCreate(ctx, req.User, items)

	if err != nil {
		log.Printf("usecase error : order create : %v, userID = %v", err.Error(), req.User)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := desc.OrderCreateResponse{
		OrderID: orderID,
	}

	return &resp, nil
}

func (s *Service) OrderPay(ctx context.Context, req *desc.OrderPayRequest) (*desc.OrderPayResponse, error) {
	err := s.usecase.OrderPay(ctx, req.OrderID)
	if err != nil {
		log.Printf("usecase error : order pay : %v for orderId = %v", err.Error(), req.OrderID)
		if errors.Is(err, dto.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, dto.ErrOrderCancelled) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		if errors.Is(err, dto.ErrOrderNotAwaitingPayment) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.OrderPayResponse{}, nil
}

func (s *Service) OrderCancel(ctx context.Context, req *desc.OrderCancelRequest) (*desc.OrderCancelResponse, error) {
	err := s.usecase.OrderCancel(ctx, req.OrderID)

	if err != nil {
		if errors.Is(err, dto.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, dto.ErrCannotCancelOrder) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.OrderCancelResponse{}, nil
}

func (s *Service) StocksInfo(ctx context.Context, req *desc.StocksInfoRequest) (*desc.StocksInfoResponse, error) {
	count, err := s.usecase.StocksInfo(ctx, req.GetSku())
	if err != nil {
		log.Printf("usecase : stocksInfo : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := desc.StocksInfoResponse{
		Count: count,
	}

	return &resp, nil
}

func convertToDtoItems(items []*desc.OrderItem) []dto.Item {
	if items == nil {
		return []dto.Item{}
	}

	result := make([]dto.Item, len(items))
	for i, item := range items {
		result[i] = dto.Item{
			SKU:   item.GetSKU(),
			Count: item.GetCount(),
		}
	}

	return result
}

func convertToOrderItems(items []dto.Item) []*desc.OrderItem {
	if items == nil {
		return []*desc.OrderItem{}
	}

	result := make([]*desc.OrderItem, len(items))
	for i, item := range items {
		result[i] = &desc.OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	return result
}
