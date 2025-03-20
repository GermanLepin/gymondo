//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=contract_mock_test.go -package=$GOPACKAGE
package service

import (
	"context"

	"gymondo/internal/model"
)

type Repository interface {
	GetProduct(ctx context.Context, productID string) (model.Product, error)
	GetProducts(ctx context.Context) ([]model.Product, error)
	GetUser(ctx context.Context, userID string) (model.User, error)
	SaveSubscription(ctx context.Context, subscription model.Subscription) error
	GetSubscription(ctx context.Context, subscriptionID string) (model.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription model.Subscription) error
	GetVoucherByCode(ctx context.Context, voucherCode string) (model.Voucher, error)
}
