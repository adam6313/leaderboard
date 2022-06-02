package v1

import (
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func (s *Server) SetRouter() {
	// middleware
	s.App.Use(
		recover.New(),
		logger.New(),
		Cros(),
	)

	// get version
	s.App.Get("/", HandleFunc(s.Version))

	r := s.App.Party("/api/v1")
	{
		// save score
		r.Post("/score", HandleFunc(s.SaveScore))

		// This endpoint can be created when the clientID is duplicated
		r.Post("/dup/score", HandleFunc(s.SaveScoreIgnoreDuplicate))

		// get LeaderBoard
		r.Get("/leaderboard", HandleFunc(s.GetLeaderBoard))
	}
}
