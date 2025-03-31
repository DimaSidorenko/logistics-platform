package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"route256/cart/internal/models"
	cartDto "route256/cart/internal/usecases/cart/dto"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestHandler_AddItem(t *testing.T) {
	mc := minimock.NewController(t)

	userID := cartDto.UserID(123)
	skuID := cartDto.SkuID(456)

	tests := []struct {
		name           string
		storageErr     error
		productErr     error
		productReturn  models.Product
		expectedErrMsg string
		before         func(*ProductClientMock, *StorageMock)
	}{
		{
			name:           "success",
			storageErr:     nil,
			productErr:     nil,
			productReturn:  models.Product{},
			expectedErrMsg: "",
			before: func(p *ProductClientMock, s *StorageMock) {
				p.GetItemMock.Expect(minimock.AnyContext, int64(skuID)).Return(models.Product{}, nil)
				s.AddItemMock.Expect(userID, skuID, uint32(2)).Return(nil)
			},
		},
		{
			name:           "validation error",
			storageErr:     nil,
			productErr:     errors.New("unknown item"),
			productReturn:  models.Product{},
			expectedErrMsg: "validate item",
			before: func(p *ProductClientMock, _ *StorageMock) {
				p.GetItemMock.Expect(minimock.AnyContext, int64(skuID)).Return(models.Product{}, errors.New("unknown item"))
			},
		},
		{
			name:           "add item error",
			storageErr:     errors.New("storage error"),
			productErr:     nil,
			productReturn:  models.Product{},
			expectedErrMsg: "add item",
			before: func(p *ProductClientMock, s *StorageMock) {
				p.GetItemMock.Expect(minimock.AnyContext, int64(skuID)).Return(models.Product{}, nil)
				s.AddItemMock.Expect(userID, skuID, uint32(2)).Return(errors.New("storage error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productClientMock := NewProductClientMock(mc)
			storageMock := NewStorageMock(mc)

			tt.before(productClientMock, storageMock)

			handler := NewHandler(productClientMock, NewLomsClientMock(mc), storageMock)
			err := handler.AddItem(userID, skuID, 2)

			if tt.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			}
		})
	}

}

func TestHandler_DeleteItem(t *testing.T) {
	mc := minimock.NewController(t)

	tests := []struct {
		name       string
		storageErr error
	}{
		{
			name:       "success",
			storageErr: nil,
		},
		{
			name:       "delete item error",
			storageErr: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewStorageMock(mc)
			storageMock.DeleteItemMock.Expect(123, 456).Return(tt.storageErr)

			handler := NewHandler(nil, nil, storageMock)
			err := handler.DeleteItem(123, 456)

			if nil == tt.storageErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	mc := minimock.NewController(t)

	tests := []struct {
		name       string
		storageErr error
	}{
		{
			name:       "success",
			storageErr: nil,
		},
		{
			name:       "delete user error",
			storageErr: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewStorageMock(mc)
			storageMock.DeleteUserMock.Expect(123).Return(tt.storageErr)

			handler := NewHandler(nil, nil, storageMock)
			err := handler.DeleteUser(123)

			if nil == tt.storageErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestHandler_GetItems(t *testing.T) {
	mc := minimock.NewController(t)

	tests := []struct {
		name           string
		storageReturn  []cartDto.Item
		storageErr     error
		productReturn  models.Product
		productErr     error
		expectedErr    error
		expectedResult cartDto.GetItemsResponse
	}{
		{
			"success",
			[]cartDto.Item{{Sku: 455, Count: 2}, {Sku: 456, Count: 3}, {Sku: 457, Count: 2}},
			nil,
			models.Product{Name: "ItemName", Price: 100},
			nil,
			nil,
			cartDto.GetItemsResponse{
				Items: []cartDto.Item{
					{Sku: 455, Count: 2, Name: "ItemName", Price: 100},
					{Sku: 456, Count: 3, Name: "ItemName", Price: 100},
					{Sku: 457, Count: 2, Name: "ItemName", Price: 100},
				},
				TotalPrice: 700,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewStorageMock(mc)
			storageMock.GetItemsMock.Expect(123).Return(tt.storageReturn, tt.storageErr)

			productClientMock := NewProductClientMock(mc)
			if nil == tt.storageErr && len(tt.storageReturn) > 0 {
				productClientMock.GetItemMock.Set(func(_ context.Context, _ int64) (models.Product, error) {
					return tt.productReturn, tt.productErr
				})
			}

			handler := NewHandler(productClientMock, nil, storageMock)
			result, err := handler.GetItems(123)

			assert.Equal(t, tt.expectedErr, err)
			assert.ElementsMatch(t, tt.expectedResult.Items, result.Items)
		})
	}
}
