package controller

import (
	"leaderboard/config"
	leaderboard_v1 "leaderboard/internal/leaderboard/interface/controller/v1"
	"leaderboard/internal/leaderboard/usecase/score"
	"net/http"

	"github.com/kataras/iris/v12"
)

// NewHTTPServer -
func NewHTTPServer(conf config.Config, scoreUsecase score.ScoreUsecase) http.Handler {
	h := leaderboard_v1.Server{
		App:          iris.New(),
		ScoreUsecase: scoreUsecase,
	}

	h.SetRouter()

	return h.App
}
