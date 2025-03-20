package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gymondo/internal/model"
	"testing"
)

func Test_Service_FindProduct(t *testing.T) {
	t.Parallel()

	t.Run("successful product fetch", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		productID := uuid.New()
		expectedProduct := model.Product{ID: productID, Name: "Test Product", Price: 100}

		mockRepo.EXPECT().GetProduct(gomock.Any(), productID.String()).Return(expectedProduct, nil)

		product, err := service.FindProduct(context.Background(), productID.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("product not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		productID := "999"

		mockRepo.EXPECT().GetProduct(gomock.Any(), productID).Return(model.Product{}, fmt.Errorf("product not found"))

		_, err := service.FindProduct(context.Background(), productID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found")
	})
}

func Test_Service_FindProducts(t *testing.T) {
	t.Parallel()

	t.Run("successful products fetch", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		expectedProducts := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 100, Tax: 10, TotalPrice: 110},
			{ID: uuid.New(), Name: "Product 2", Price: 200, Tax: 20, TotalPrice: 220},
		}

		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(expectedProducts, nil)

		products, err := service.FindProducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("error fetching products", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		_, err := service.FindProducts(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func Test_Service_FindProductsWithVoucher(t *testing.T) {
	t.Parallel()

	t.Run("successful products didn't find", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "DISCOUNT10"
		voucher := model.Voucher{
			Code:          voucherCode,
			DiscountType:  model.Percentage,
			DiscountValue: 0.1,
		}

		var products []model.Product

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(voucher, nil)
		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(products, nil)

		_, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.NoError(t, err)
	})

	t.Run("successful products fetch with voucher (percentage)", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "DISCOUNT10"
		voucher := model.Voucher{
			Code:          voucherCode,
			DiscountType:  model.Percentage,
			DiscountValue: 0.1,
		}

		products := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 100, Tax: 10, TotalPrice: 110},
			{ID: uuid.New(), Name: "Product 2", Price: 200, Tax: 20, TotalPrice: 220},
		}

		expectedProducts := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 90, Tax: 9, TotalPrice: 99},
			{ID: uuid.New(), Name: "Product 2", Price: 180, Tax: 18, TotalPrice: 198},
		}

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(voucher, nil)
		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(products, nil)

		resultProducts, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.NoError(t, err)
		assert.Len(t, resultProducts, len(expectedProducts))

		for i := range expectedProducts {
			assert.Equal(t, expectedProducts[i].Price, resultProducts[i].Price)
			assert.Equal(t, expectedProducts[i].Tax, resultProducts[i].Tax)
			assert.Equal(t, expectedProducts[i].TotalPrice, resultProducts[i].TotalPrice)
		}
	})

	t.Run("successful products fetch with voucher (fixed)", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "Fixed15"
		voucher := model.Voucher{
			Code:          voucherCode,
			DiscountType:  model.Fixed,
			DiscountValue: 20,
		}

		products := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 90, Tax: 10, TotalPrice: 100},
			{ID: uuid.New(), Name: "Product 2", Price: 200, Tax: 20, TotalPrice: 220},
		}

		expectedProducts := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 72, Tax: 8, TotalPrice: 80},
			{ID: uuid.New(), Name: "Product 2", Price: 181.81, Tax: 18.19, TotalPrice: 200},
		}

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(voucher, nil)
		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(products, nil)

		resultProducts, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.NoError(t, err)
		assert.Len(t, resultProducts, len(expectedProducts))

		for i := range expectedProducts {
			assert.Equal(t, expectedProducts[i].Price, resultProducts[i].Price)
			assert.Equal(t, expectedProducts[i].Tax, resultProducts[i].Tax)
			assert.Equal(t, expectedProducts[i].TotalPrice, resultProducts[i].TotalPrice)
		}
	})

	t.Run("failed to calculate products (fixed)", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "Fixed15"
		voucher := model.Voucher{
			Code:          voucherCode,
			DiscountType:  model.Fixed,
			DiscountValue: 620,
		}

		products := []model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 90, Tax: 10, TotalPrice: 100},
			{ID: uuid.New(), Name: "Product 2", Price: 200, Tax: 20, TotalPrice: 220},
		}

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(voucher, nil)
		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(products, nil)

		expectedError := errors.New("product price couldn't be less than 0")
		_, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("voucher not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "INVALID"

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(model.Voucher{}, fmt.Errorf("voucher not found"))

		_, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch voucher")
	})

	t.Run("error fetching products", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		voucherCode := "DISCOUNT10"
		voucher := model.Voucher{
			Code:          voucherCode,
			DiscountType:  model.Percentage,
			DiscountValue: 0.1,
		}

		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(voucher, nil)
		mockRepo.EXPECT().GetProducts(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		_, err := service.FindProductsWithVoucher(context.Background(), voucherCode)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch products")
	})
}
