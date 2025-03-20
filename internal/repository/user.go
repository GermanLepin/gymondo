package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gymondo/internal/model"
)

func (r *Repository) GetUser(
	ctx context.Context,
	userID string,
) (model.User, error) {
	const query = `
		select id, first_name, second_name, email
		from service.users
		where id = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.SecondName,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user with ID %s not found: %w", userID, err)
		}
		return user, fmt.Errorf("failed to retrieve user with ID %s: %w", userID, err)
	}

	return user, nil
}
