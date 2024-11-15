package repository

import (
	"context"
	"time"

	"github.com/whitexwc/basic-go/webook/internal/domain"
	"github.com/whitexwc/basic-go/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	// 在这里操作缓存
}

func (r *UserRepository) FindById(int64) {
	// 先从cache里面找
	// 再从dao里面找
	// 找到了回写cache
}

func (r *UserRepository) UpdateProfileById(ctx context.Context, u domain.User) error {
	return r.dao.UpdateByUserId(ctx, dao.User{
		Id:       u.Id,
		Birthday: u.Birthday.UnixMilli(),
		NickName: u.NickName,
		AboutMe:  u.AboutMe,
	})
}

func (r *UserRepository) FindProfileById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, id)
	return domain.User{
		NickName: u.NickName,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}, err
}
