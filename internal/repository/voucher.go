package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gymondo/internal/model"
)

func (r *Repository) GetVoucherByCode(
	ctx context.Context,
	voucherCode string,
) (model.Voucher, error) {
	const query = `
		select id, code, discount_type, discount_value
		from service.vouchers
		where code = $1
	`

	var voucher model.Voucher
	err := r.db.QueryRowContext(ctx, query, voucherCode).Scan(
		&voucher.ID,
		&voucher.Code,
		&voucher.DiscountType,
		&voucher.DiscountValue,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return voucher, fmt.Errorf("voucher with code %s not found", voucherCode)
		}
		return voucher, fmt.Errorf("failed to query voucher with code %s: %w", voucherCode, err)
	}

	return voucher, nil
}
