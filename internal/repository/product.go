package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gymondo/internal/model"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetProducts(ctx context.Context) ([]model.Product, error) {
	const query = `
		select id, name, duration_days, price, tax, total_price
		from service.products
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.DurationDays,
			&product.Price,
			&product.Tax,
			&product.TotalPrice,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over products: %w", err)
	}

	return products, nil
}

func (r *Repository) GetProduct(
	ctx context.Context,
	productID string,
) (model.Product, error) {
	const query = `
		select id, name, duration_days, price, tax, total_price
		from service.products
		where id = $1
	`

	var product model.Product
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.DurationDays,
		&product.Price,
		&product.Tax,
		&product.TotalPrice,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, fmt.Errorf("product with ID %s not found: %w", productID, err)
		}
		return product, fmt.Errorf("failed to query product by ID %s: %w", productID, err)
	}

	return product, nil
}
