package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gymondo/internal/model"
)

func (s *Service) Subscribe(
	ctx context.Context,
	userID string,
	productID string,
	voucherCode string,
	trialPeriod bool,
) (string, error) {
	user, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	product, err := s.repository.GetProduct(ctx, productID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch product: %w", err)
	}

	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, product.DurationDays).Truncate(24 * time.Hour)

	subscriptionID := uuid.New()
	subscription := model.Subscription{
		ID:           subscriptionID,
		UserID:       user.ID,
		ProductID:    product.ID,
		StartDate:    startDate,
		EndDate:      endDate,
		DurationDays: product.DurationDays,
		Price:        product.Price,
		Tax:          product.Tax,
		TotalPrice:   product.TotalPrice,
		Status:       model.Active,
	}

	if voucherCode != "" {
		voucher, err := s.repository.GetVoucherByCode(ctx, voucherCode)
		if err != nil {
			return "", fmt.Errorf("failed to fetch voucher: %w", err)
		}

		productWithVoucher, err := calculatePriceWithVoucher(product, voucher)
		if err != nil {
			return "", fmt.Errorf("failed to calculate price with voucher: %w", err)
		}

		subscription.Price = productWithVoucher.Price
		subscription.Tax = productWithVoucher.Tax
		subscription.TotalPrice = productWithVoucher.TotalPrice
	}
	if trialPeriod {
		subscription.TrialStartDate = &startDate
		trialEndDate := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
		subscription.TrialEndDate = &trialEndDate
	}

	if err := s.repository.SaveSubscription(ctx, subscription); err != nil {
		return "", fmt.Errorf("failed to save subscription: %w", err)
	}

	return subscriptionID.String(), nil
}

func (s *Service) FindSubscription(ctx context.Context, subscriptionID string) (model.Subscription, error) {
	subscription, err := s.repository.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("failed to fetch subscription with ID %s: %w", subscriptionID, err)
	}

	return subscription, nil
}

func (s *Service) PauseSubscription(ctx context.Context, subscriptionID string) error {
	subscription, err := s.repository.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to find subscription with ID %s: %w", subscriptionID, err)
	}

	switch subscription.Status {
	case model.Paused:
		return fmt.Errorf("subscription is already paused")
	case model.Canceled:
		return fmt.Errorf("subscription is canceled")
	}

	if subscription.TrialEndDate != nil {
		if subscription.TrialEndDate.After(time.Now().Truncate(24 * time.Hour)) {
			return fmt.Errorf("can't pause subscription during trial period")
		}
	}

	subscription.Status = model.Paused
	pausedDate := time.Now().Truncate(24 * time.Hour)
	subscription.PausedDate = &pausedDate

	err = s.repository.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to pause subscription: %w", err)
	}

	return nil
}

func (s *Service) UnpauseSubscription(ctx context.Context, subscriptionID string) error {
	subscription, err := s.repository.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to find subscription with ID %s: %w", subscriptionID, err)
	}

	switch subscription.Status {
	case model.Active:
		return fmt.Errorf("subscription is already active")
	case model.Canceled:
		return fmt.Errorf("subscription is canceled")
	}

	subscription.Status = model.Active
	unpausedDate := time.Now().Truncate(24 * time.Hour)
	subscription.UnpausedDate = &unpausedDate

	err = s.repository.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to pause subscription: %w", err)
	}

	return nil
}

func (s *Service) CancelSubscription(ctx context.Context, subscriptionID string) error {
	subscription, err := s.repository.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to find subscription with ID %s: %w", subscriptionID, err)
	}

	switch subscription.Status {
	case model.Paused:
		return fmt.Errorf("subscription is paused")
	case model.Canceled:
		return fmt.Errorf("subscription is already canceled")
	}

	subscription.Status = model.Canceled
	canceledDate := time.Now().Truncate(24 * time.Hour)
	subscription.CanceledDate = &canceledDate

	err = s.repository.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to pause subscription: %w", err)
	}

	return nil
}
