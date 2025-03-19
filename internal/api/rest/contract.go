package rest

import (
	"context"

	"gymondo/internal/model"
)

type productService interface {
	FindOne(ctx context.Context, productID string) (model.Product, error)
	FindAll(ctx context.Context) ([]model.Product, error)
}
