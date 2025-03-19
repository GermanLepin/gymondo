package model

import (
	"time"
	
	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	Active   SubscriptionStatus = "active"
	Paused   SubscriptionStatus = "paused"
	Canceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID                    uuid.UUID          `json:"id"`
	UserID                uuid.UUID          `json:"user_id"`
	ProductID             uuid.UUID          `json:"product_id"`
	StartDate             time.Time          `json:"start_date"`
	EndDate               time.Time          `json:"end_date"`
	DurationDays          int                `json:"duration_days"`
	Price                 int64              `json:"price"`
	Tax                   int64              `json:"tax"`
	TotalPrice            int64              `json:"total_price"`
	TotalPriceWithVoucher int64              `json:"total_price_with_voucher"`
	Status                SubscriptionStatus `json:"status"`
	TrialStartDate        *time.Time         `json:"trial_start_date,omitempty"`
	TrialEndDate          *time.Time         `json:"trial_end_date,omitempty"`
	CanceledDate          *time.Time         `json:"canceled_date,omitempty"`
	PausedDate            *time.Time         `json:"paused_date,omitempty"`
}
