package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, u *User) error
	Store(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, u *User) error
	Store(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
}
