package memory

import (
	"context"
	"leaderboard/internal/leaderboard/domain/model"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

// Create -
func (r *Repo) Create(ctx context.Context, key string, in *model.Score) error {
	return r.client.ZAdd(ctx, key, &goredis.Z{
		Score:  float64(in.Score),
		Member: in.ClientID,
	}).Err()
}

// List
func (r *Repo) List(ctx context.Context, key string, offset, limit int64) ([]*model.Score, error) {
	scores, err := r.client.ZRevRangeWithScores(ctx, key, offset, limit).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*model.Score, len(scores))

	for i, z := range scores {
		result[i] = &model.Score{
			ClientID: z.Member.(string),
			Score:    z.Score,
		}
	}

	return result, nil
}

// DeleteAll
func (r *Repo) DeleteAll(ctx context.Context) error {
	iter := r.client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		r.client.Del(ctx, iter.Val())
	}

	return nil
}

// SetExpire set key expire(TTL)
func (r *Repo) SetExpire(ctx context.Context, key string, t time.Duration) error {
	_, err := r.client.Expire(ctx, key, time.Minute*10).Result()
	if err != nil {
		return err
	}

	return err
}

// Exists check key is exist
func (r *Repo) Exists(ctx context.Context, key string) int64 {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return 0
	}

	return count
}
