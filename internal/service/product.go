package service

import (
	"context"

	"gymondo/internal/model"
)

type Service struct {
	productRepository productRepository
}

func New(productRepository productRepository) *Service {
	return &Service{
		productRepository: productRepository,
	}
}

func (s *Service) FindOne(ctx context.Context, productID string) (model.Product, error) {
	return s.productRepository.Get(ctx, productID)
}

func (s *Service) FindAll(ctx context.Context) ([]model.Product, error) {
	return s.productRepository.GetList(ctx)
}
