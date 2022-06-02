package repository

import (
	"context"
	"leaderboard/internal/leaderboard/domain/model"
	"time"
)

// LeaderBoardRepository Repository interface for leaderboard service
type LeaderBoardRepository interface {
	// Create
	Create(ctx context.Context, key string, score *model.Score) error

	// List
	List(ctx context.Context, key string, offset, limit int64) ([]*model.Score, error)

	// DeleteAll
	DeleteAll(ctx context.Context) error

	// SetExpire
	SetExpire(ctx context.Context, key string, t time.Duration) error

	// Exists
	Exists(ctx context.Context, key string) int64
}
