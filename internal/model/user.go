package model

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Email      string    `json:"email"`
}
