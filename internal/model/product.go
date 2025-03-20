package model

import "github.com/google/uuid"

type Product struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	DurationDays int       `json:"duration_days"`
	Price        float64   `json:"price"`
	Tax          float64   `json:"tax"`
	TotalPrice   float64   `json:"total_price"`
}
