package model

import (
	"time"
	
	"github.com/google/uuid"
)

type VoucherStatus string

const (
	Percentage VoucherStatus = "percentage"
	Fixed      VoucherStatus = "fixed"
)

// Voucher represents a discount voucher
type Voucher struct {
	ID            uuid.UUID     `json:"id"`
	Code          string        `json:"code"`
	DiscountType  VoucherStatus `json:"discount_type"`
	DiscountValue int64         `json:"discount_value"`
	ValidFrom     time.Time     `json:"valid_from"`
	ValidUntil    time.Time     `json:"valid_until"`
}
