package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/whitexwc/basic-go/webook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 传单机 Redis 也可以传cluster的Redis
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到了 B， B一定是接口 => 保证面向接口编程
// A 用到了 B，B 一定是A的字段 => 规避包变量，包方法(缺乏扩展性)
// A 用到了 B， A绝对不初始化B，而是外面注入 => 保持依赖注入和依赖反转
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 只要 error 为 nil， 就认为缓存有山上
// 如果没有数据，返回一个特定的error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	// 数据不存在，返回的是redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
