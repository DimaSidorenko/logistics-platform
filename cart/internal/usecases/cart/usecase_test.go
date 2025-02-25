package cart

import (
	"errors"
	"fmt"
	cartDto "route256/cart/internal/usecases/cart/dto"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"route256/cart/internal/models"
)

func TestHandler_AddItem(t *testing.T) {
	mc := minimock.NewController(t)

	tests := []struct {
		name           string
		storageErr     error
		productErr     error
		productReturn  models.Product
		expectedErrMsg string
	}{
		{
			name:           "success",
			storageErr:     nil,
			productErr:     nil,
			productReturn:  models.Product{},
			expectedErrMsg: "",
		},
		{
			name:           "validation error",
			storageErr:     nil,
			productErr:     errors.New("unknown item"),
			productReturn:  models.Product{},
			expectedErrMsg: "validate item",
		},
		{
			name:           "add item error",
			storageErr:     errors.New("storage error"),
			productErr:     nil,
			productReturn:  models.Product{},
			expectedErrMsg: "add item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productClientMock := NewProductClientMock(mc)
			productClientMock.GetItemMock.Expect(int64(456)).Return(tt.productReturn, tt.productErr)

			storageMock := NewStorageMock(mc)
			if nil == tt.productErr {
				storageMock.AddItemMock.Expect(123, 456, uint32(2)).Return(tt.storageErr)
			}

			handler := NewHandler(productClientMock, storageMock)
			err := handler.AddItem(123, 456, 2)

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

			handler := NewHandler(nil, storageMock)
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

			handler := NewHandler(nil, storageMock)
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
			[]cartDto.Item{{Sku: 456, Count: 2}},
			nil,
			models.Product{Name: "ItemName", Price: 100},
			nil,
			nil,
			cartDto.GetItemsResponse{
				Items:      []cartDto.Item{{Sku: 456, Count: 2, Name: "ItemName", Price: 100}},
				TotalPrice: 200,
			},
		},
		{
			"storage error",
			nil,
			errors.New("storage error"),
			models.Product{},
			nil,
			errors.New("get items: storage error"),
			cartDto.GetItemsResponse{},
		},
		{
			"user not found",
			[]cartDto.Item{},
			nil,
			models.Product{},
			nil,
			cartDto.ErrUserNotFound,
			cartDto.GetItemsResponse{},
		},
		{
			"item not found",
			[]cartDto.Item{{Sku: 456, Count: 2}},
			nil,
			models.Product{},
			models.ErrItemNotFound,
			fmt.Errorf("get item in product service: %v", models.ErrItemNotFound),
			cartDto.GetItemsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewStorageMock(mc)
			storageMock.GetItemsMock.Expect(123).Return(tt.storageReturn, tt.storageErr)

			productClientMock := NewProductClientMock(mc)
			if nil == tt.storageErr && len(tt.storageReturn) > 0 {
				productClientMock.GetItemMock.Expect(int64(456)).Return(tt.productReturn, tt.productErr)
			}

			handler := NewHandler(productClientMock, storageMock)
			result, err := handler.GetItems(123)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
