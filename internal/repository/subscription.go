package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gymondo/internal/model"
	"time"
)

func (r *Repository) SaveSubscription(ctx context.Context, subscription model.Subscription) error {
	query := `
		INSERT INTO service.subscriptions (
			id, 
		    user_id, 
		    product_id, 
		    start_date,
		    end_date,
		    duration_days, 
			price, 
		    tax, 
		    total_price, 
		    status, 
			trial_start_date, 
		    trial_end_date, 
		    canceled_date,
		    paused_date,
			unpaused_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.UserID,
		subscription.ProductID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.DurationDays,
		subscription.Price,
		subscription.Tax,
		subscription.TotalPrice,
		subscription.Status,
		nullTime(subscription.TrialStartDate),
		nullTime(subscription.TrialEndDate),
		nullTime(subscription.CanceledDate),
		nullTime(subscription.PausedDate),
		nullTime(subscription.UnpausedDate),
	)
	if err != nil {
		return fmt.Errorf("failed to save subscription with ID %s: %w", subscription.ID, err)
	}

	return nil
}

func nullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func (r *Repository) GetSubscription(ctx context.Context, subscriptionID string) (model.Subscription, error) {
	query := `
		select 
			id, 
			user_id, 
			product_id, 
			start_date,
			end_date,
			duration_days, 
			price, 
			tax, 
			total_price,
			status, 
			trial_start_date, 
			trial_end_date, 
			canceled_date,
			paused_date,
			unpaused_date
		from service.subscriptions 
		where id = $1
	`

	var subscription model.Subscription
	err := r.db.QueryRowContext(ctx, query, subscriptionID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.ProductID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.DurationDays,
		&subscription.Price,
		&subscription.Tax,
		&subscription.TotalPrice,
		&subscription.Status,
		&subscription.TrialStartDate,
		&subscription.TrialEndDate,
		&subscription.CanceledDate,
		&subscription.PausedDate,
		&subscription.UnpausedDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subscription, fmt.Errorf("subscription with ID %s not found: %w", subscriptionID, err)
		}
		return subscription, fmt.Errorf("failed to retrieve subscription with ID %s: %w", subscriptionID, err)
	}

	return subscription, nil
}

func (r *Repository) UpdateSubscription(
	ctx context.Context,
	subscription model.Subscription,
) error {
	query := `
		UPDATE service.subscriptions
		SET 
		    status = $2,
			canceled_date = $3,
			paused_date = $4,
			unpaused_date = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.Status,
		subscription.CanceledDate,
		subscription.PausedDate,
		subscription.UnpausedDate,
	)
	if err != nil {
		return fmt.Errorf("failed to update subscription with ID %s: %w", subscription.ID, err)
	}

	return nil
}
