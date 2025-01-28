package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewUserStorage(rdb *redis.Client) UserStorage {
	return &userStore{
		rdb: rdb,
	}
}

type userStore struct {
	rdb *redis.Client
}

type UserStorage interface{}
