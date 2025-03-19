//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=contract_mock_test.go -package=$GOPACKAGE
package service

import (
	"context"

	"gymondo/internal/model"
)

type productRepository interface {
	Get(ctx context.Context, productID string) (model.Product, error)
	GetList(ctx context.Context) ([]model.Product, error)
}
