package user

import (
	"context"

	"github.com/budhilaw/blog/models"
)

// Usecase for business logic
type Usecase interface {
	Fetch(ctx context.Context, num int64) (res []*models.User, err error)
	Store(ctx context.Context, a *models.User) error
	Login(ctx context.Context, email string, password string, a *models.User) error
}
