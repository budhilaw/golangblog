package usecase

import (
	"context"
	"time"

	"github.com/budhilaw/blog/models"

	"github.com/budhilaw/blog/user"
)

type userUsecase struct {
	//postRepo post.Repository
	userRepo       user.Repository
	contextTimeout time.Duration
}

// New : Will create object that represent of the user.Usecase interface
func New(u user.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Fetch(c context.Context, num int64) ([]*models.User, error) {
	if num == 0 {
		num = 10
	}

	ctx, ctxCancel := context.WithTimeout(c, u.contextTimeout)
	defer ctxCancel()

	listUser, err := u.userRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	return listUser, nil
}

func (u *userUsecase) Store(c context.Context, m *models.User) error {
	ctx, ctxCancel := context.WithTimeout(c, u.contextTimeout)
	defer ctxCancel()

	err := u.userRepo.Store(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Login(c context.Context, email string, pass string, m *models.User) error {
	ctx, ctxCancel := context.WithTimeout(c, u.contextTimeout)
	defer ctxCancel()

	err := u.userRepo.Login(ctx, email, pass, m)
	if err != nil {
		return err
	}

	return nil
}
