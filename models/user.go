package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreadedAt time.Time `json:"created_at"`
}
