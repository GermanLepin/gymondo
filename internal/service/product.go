package service

import (
	"context"
	"fmt"
	"gymondo/internal/model"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindProduct(ctx context.Context, productID string) (model.Product, error) {
	return s.repository.GetProduct(ctx, productID)
}

func (s *Service) FindProducts(ctx context.Context) ([]model.Product, error) {
	return s.repository.GetProducts(ctx)
}

func (s *Service) FindProductsWithVoucher(ctx context.Context, voucherCode string) ([]model.Product, error) {
	voucher, err := s.repository.GetVoucherByCode(ctx, voucherCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch voucher %s: %w", voucherCode, err)
	}

	products, err := s.repository.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	if len(products) == 0 {
		return []model.Product{}, nil
	}

	responseProducts := make([]model.Product, 0, len(products))
	for _, product := range products {
		productWithVoucher, err := calculatePriceWithVoucher(product, voucher)
		if err != nil {
			return []model.Product{}, fmt.Errorf("failed to calculate products: %w", err)
		}
		responseProducts = append(responseProducts, model.Product{
			ID:         product.ID,
			Name:       product.Name,
			Price:      productWithVoucher.Price,
			Tax:        productWithVoucher.Tax,
			TotalPrice: productWithVoucher.TotalPrice,
		})
	}

	return responseProducts, nil
}
