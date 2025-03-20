package model

import (
	"github.com/google/uuid"
)

type VoucherStatus string

const (
	Percentage VoucherStatus = "percentage"
	Fixed      VoucherStatus = "fixed"
)

type Voucher struct {
	ID            uuid.UUID     `json:"id"`
	Code          string        `json:"code"`
	DiscountType  VoucherStatus `json:"discount_type"`
	DiscountValue float64       `json:"discount_value"`
}
