package v1

import (
	"errors"
	"leaderboard/config"
	"leaderboard/internal/leaderboard/usecase/score"

	"github.com/kataras/iris/v12"
)

type Server struct {
	App          *iris.Application
	ScoreUsecase score.ScoreUsecase
}

// Version used to get version, and ping pong check
func (s *Server) Version(c *C) {
	c.R(map[string]string{
		"version": config.C.Version,
	})
}

// SaveScore -
func (s *Server) SaveScore(c *C) {
	// get clientId from head
	clientId := c.Request().Header.Get("ClientId")

	// check clientId is exist
	if clientId == "" {
		c.E(errors.New("bad request"))
		return
	}

	// get body data
	data := &score.AddScore{
		ClientID: clientId,
	}
	if err := c.ReadJSON(data); err != nil {
		c.E(err)
		return
	}

	// usecase
	if err := s.ScoreUsecase.Add(c.Request().Context(), data); err != nil {
		c.E(err)
		return
	}

	c.R(nil)
}

// SaveScoreIgnoreDuplicate
func (s *Server) SaveScoreIgnoreDuplicate(c *C) {
	// get clientId from head
	clientId := c.Request().Header.Get("ClientId")

	// check clientId is exist
	if clientId == "" {
		c.E(errors.New("bad request"))
		return
	}

	// get body data
	data := &score.AddScore{
		ClientID: clientId,
	}
	if err := c.ReadJSON(data); err != nil {
		c.E(err)
		return
	}

	// usecase
	if err := s.ScoreUsecase.AddIgnoreDuplicate(c.Request().Context(), data); err != nil {
		c.E(err)
		return
	}

	c.R(nil)
}

// GetLeaderBoard
func (s *Server) GetLeaderBoard(c *C) {
	scores, err := s.ScoreUsecase.GetLeaderBoard(c.Request().Context())
	if err != nil {
		c.E(err)
		return
	}

	c.R(map[string]interface{}{
		"topPlayers": scores,
	})
}
