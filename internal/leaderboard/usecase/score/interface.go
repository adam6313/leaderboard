package score

import (
	"context"
	"leaderboard/internal/leaderboard/domain/model"
)

// ScoreUsecase -
type ScoreUsecase interface {
	// Add - add score
	Add(ctx context.Context, command *AddScore) error

	// AddIgnoreDuplicate
	AddIgnoreDuplicate(ctx context.Context, command *AddScore) error

	// GetLeaderBoard
	GetLeaderBoard(ctx context.Context) ([]*model.Score, error)

	// ResetLeaderBoard
	ResetLeaderBoard(ctx context.Context) error
}
