package cache

import "github.com/redis/go-redis/v9"

type Storage struct {
	Users UserStorage
}

func NewRedisStorage(rdb *redis.Client) *Storage {
	return &Storage{
		Users: NewUserStorage(rdb),
	}
}
