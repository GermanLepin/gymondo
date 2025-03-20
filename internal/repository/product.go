package product

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

func NewProductRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(
	ctx context.Context,
	productID string,
) (model.Product, error) {
	const query = `
		select id, name, price, duration_days
		from service.products
		where id = $1
	`

	var product model.Product
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.DurationDays,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, fmt.Errorf("product with ID %s not found", productID)
		}
		return product, err
	}

	return product, nil
}

func (r *Repository) GetList(ctx context.Context) ([]model.Product, error) {
	const query = `
		SELECT id, name, price, duration_days
		FROM service.products
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.DurationDays,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
