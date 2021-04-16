package usecase

import (
	"context"
	"time"

	"github.com/diantanjung/blogo/user-service/domain"
	validator "gopkg.in/go-playground/validator.v9"
)

type userUsecase struct {
	userRepo domain.UserRepository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an userUsecase object representation of domain.userUsecase interface
func NewUserUsecase(a domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: a,
	}
}
func (a *userUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.userRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

func (a *userUsecase) GetByID(c context.Context, id int64) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *userUsecase) Update(c context.Context, u *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	u.UpdatedAt = time.Now()

	if err = isUserValid(u); err != nil {
		return
	}
	return a.userRepo.Update(ctx, u)
}
func (a *userUsecase) Store(c context.Context,u *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	if err = isUserValid(u); err != nil {
		return
	}
	return a.userRepo.Store(ctx, u)
}
func (a *userUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.User{}) {
		return domain.ErrNotFound
	}
	return a.userRepo.Delete(ctx, id)
}

func isUserValid(m *domain.User) error {
	validate := validator.New()
	return validate.Struct(m)
}