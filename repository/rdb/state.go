package rdb

import (
	"context"
	"encoding/json"
	"strconv"
	"taskbot/repository"
	"taskbot/service/telegram"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	dataExpiration   = time.Hour * 2 * 24
	operationTimeout = time.Second * 10
)

type StateStore struct {
	r *redis.Client
}

func NewStateStore(r *redis.Client) *StateStore {
	return &StateStore{
		r: r,
	}
}

func (s *StateStore) Get(id int64) (telegram.ChatState, error) {
	chatState := telegram.ChatState{}

	key := strconv.FormatInt(id, 10)

	ctx, cl := context.WithTimeout(context.Background(), operationTimeout)
	defer cl()
	val, err := s.r.Get(ctx, key).Result()

	if err == redis.Nil {
		return chatState, repository.ErrNotFound
	}
	if err != nil {
		return chatState, err
	}

	err = json.Unmarshal([]byte(val), &chatState)
	return chatState, err
}

func (s *StateStore) Save(state telegram.ChatState) error {
	key := strconv.FormatInt(state.Id, 10)

	serializedBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	ctx, cl := context.WithTimeout(context.Background(), operationTimeout)
	defer cl()

	err = s.r.Set(ctx, key, string(serializedBytes), dataExpiration).Err()

	return err
}

func (s *StateStore) Delete(state telegram.ChatState) error {
	key := strconv.FormatInt(state.Id, 10)

	ctx, cl := context.WithTimeout(context.Background(), operationTimeout)
	defer cl()

	err := s.r.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}

	return err
}
