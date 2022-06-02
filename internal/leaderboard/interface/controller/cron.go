package controller

import (
	"context"
	"leaderboard/internal/leaderboard/usecase/score"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// NewCron
func NewCron(ctx context.Context, usecase score.ScoreUsecase, logger *zap.Logger) *cron.Cron {
	c := cron.New()

	c.AddFunc("*/10 * * * *", func() {
		logger.Sugar().Info("start cron job")
		retry(3, time.Duration(time.Second), usecase.ResetLeaderBoard)
	})

	return c
}

func retry(attempts int, sleep time.Duration, f func(ctx context.Context) error) error {
	if err := f(context.Background()); err != nil {
		if attempts--; attempts > 0 {
			// add some randomness
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, sleep, f)
		}
		return err
	}

	return nil
}
