//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=contract_mock_test.go -package=$GOPACKAGE
package rest

import (
	"context"

	"gymondo/internal/model"
)

type service interface {
	FindProduct(ctx context.Context, productID string) (model.Product, error)
	FindProducts(ctx context.Context) ([]model.Product, error)
	FindProductsWithVoucher(ctx context.Context, voucherCode string) ([]model.Product, error)
	Subscribe(
		ctx context.Context,
		userID string,
		productID string,
		voucherCode string,
		trialPeriod bool,
	) (subscriptionID string, err error)
	FindSubscription(ctx context.Context, subscriptionID string) (model.Subscription, error)
	PauseSubscription(ctx context.Context, subscriptionID string) error
	UnpauseSubscription(ctx context.Context, subscriptionID string) error
	CancelSubscription(ctx context.Context, subscriptionID string) error
}
