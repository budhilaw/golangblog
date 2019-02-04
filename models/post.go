package models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    User      `json:"user"`
	UpdatedAt time.Time `json:"updated_at"`
	CreadedAt time.Time `json:"created_at"`
}
