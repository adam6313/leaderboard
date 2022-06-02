package memory

import (
	"leaderboard/config"
	"leaderboard/internal/leaderboard/domain/repository"

	goredis "github.com/go-redis/redis/v8"
)

type Repo struct {
	client *goredis.Client
}

// NewRepository -
func NewRepository(client *goredis.Client, c config.Config) repository.LeaderBoardRepository {
	return &Repo{
		client: client,
	}
}
