package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ariefro/threads-server/internal/store"
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

type UserStorage interface {
	Get(context.Context, int64) (*store.User, error)
	Set(context.Context, *store.User) error
}

const UserExpTime = time.Minute

func (s *userStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *userStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}
