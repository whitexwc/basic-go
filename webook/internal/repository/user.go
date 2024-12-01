package repository

import (
	"context"
	"github.com/whitexwc/basic-go/webook/internal/repository/cache"
	"time"

	"github.com/whitexwc/basic-go/webook/internal/domain"
	"github.com/whitexwc/basic-go/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
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

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从cache里面找
	// 再从dao里面找
	// 找到了回写cache
	u, err := r.cache.Get(ctx, id)
	// 几种情况：
	// 缓存里面有数据
	// 缓存里面没数据
	// 缓存出错了，你也不知道有没有数据
	if err == nil {
		//必然有数据
		return u, nil
	}
	// 没有这个数据
	if err == cache.ErrKeyNotExist {
		// 去数据库里面加载

	}
	// 如果出现其他 error， 是否应该去数据库加载
	// 比如 err == io.EOF
	// 如果是偶发性的故障，应该选择加载
	// 如果是redis崩了，如果选择加载，需要保护好系统，有兜底，避免数据库被打崩
	// 通过数据库限流、或者布尔过滤器的方式保护数据库
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	go func() {
		err = r.cache.Set(ctx, u)
		// 缓存设置失败
		if err != nil {
			// 打日志，做监控
		}
	}()
	return u, err
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
