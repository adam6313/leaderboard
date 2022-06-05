package v1

import (
	"errors"
	"leaderboard/internal/leaderboard/domain/model"
	"leaderboard/internal/leaderboard/usecase/score"
	socre "leaderboard/test/mock/usecase"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/golang/mock/gomock"
	"github.com/kataras/iris/v12/httptest"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockScoreUsecase *socre.MockScoreUsecase
	server           *Server
	mockHTTP         *httpexpect.Expect
}

// SetupTest
func (t *handlerSuite) SetupSuite() {
	t.ctrl = gomock.NewController(t.T())

	t.mockScoreUsecase = socre.NewMockScoreUsecase(t.ctrl)

	t.server = &Server{
		App:          iris.New(),
		ScoreUsecase: t.mockScoreUsecase,
	}

	t.server.SetRouter()
	t.mockHTTP = httptest.New(t.T(), t.server.App, httptest.URL("http://localhost:8080"))
}

// TestHandler
func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

// Test_Version
func (h *handlerSuite) Test_Version() {
	h.mockHTTP.GET("/").
		Expect().
		Status(httptest.StatusOK).
		JSON().Object().
		ContainsKey("version").ValueEqual("version", "")
}

// Test_SaveScore
func (h *handlerSuite) Test_SaveScore() {
	type args struct {
		headers map[string]string
		body    interface{}
	}

	tests := []struct {
		name string
		args args
		fn   func(args) *httpexpect.Object
		want map[string]interface{}
	}{
		{
			name: "test not set Headers",
			args: args{},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/score").
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "bad request",
			},
		},
		{
			name: "test not set body",
			args: args{
				headers: map[string]string{
					"ClientId": "adam",
				},
				body: map[string]interface{}{},
			},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/score").
					WithHeaders(in.headers).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "unexpected end of JSON input",
			},
		},
		{
			name: "test wrong structure case",
			args: args{
				headers: map[string]string{
					"ClientId": "adam",
				},
				body: `{"test": 1}`,
			},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/score").
					WithHeaders(in.headers).
					WithJSON(in.body).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "json: cannot unmarshal string into Go value of type score.AddScore",
			},
		},
		{
			name: "test save score success",
			args: args{
				headers: map[string]string{
					"ClientId": "peter",
				},
				body: &score.AddScore{
					Score: 100.2,
				},
			},
			fn: func(in args) *httpexpect.Object {
				var e error
				command := in.body.(*score.AddScore)
				command.ClientID = "peter"

				h.mockScoreUsecase.EXPECT().Add(gomock.Any(), command).Return(e).Times(1)

				return h.mockHTTP.POST("/api/v1/score").
					WithHeaders(in.headers).
					WithJSON(in.body).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status")
			},
			want: map[string]interface{}{
				"status": "ok",
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {

			expect := test.fn(test.args)
			for k, w := range test.want {
				expect.ValueEqual(k, w)
			}
		})
	}
}

// Test_SaveScoreIgnoreDuplicate
func (h *handlerSuite) Test_SaveScoreIgnoreDuplicate() {
	type args struct {
		headers map[string]string
		body    interface{}
	}

	tests := []struct {
		name string
		args args
		fn   func(args) *httpexpect.Object
		want map[string]interface{}
	}{
		{
			name: "test not set Headers",
			args: args{},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/dup/score").
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "bad request",
			},
		},
		{
			name: "test not set body",
			args: args{
				headers: map[string]string{
					"ClientId": "adam",
				},
				body: map[string]interface{}{},
			},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/dup/score").
					WithHeaders(in.headers).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "unexpected end of JSON input",
			},
		},
		{
			name: "test wrong structure case",
			args: args{
				headers: map[string]string{
					"ClientId": "adam",
				},
				body: `{"test": 1}`,
			},
			fn: func(in args) *httpexpect.Object {
				return h.mockHTTP.POST("/api/v1/dup/score").
					WithHeaders(in.headers).
					WithJSON(in.body).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object().
					ContainsKey("message")
			},
			want: map[string]interface{}{
				"message": "json: cannot unmarshal string into Go value of type score.AddScore",
			},
		},
		{
			name: "test save score success",
			args: args{
				headers: map[string]string{
					"ClientId": "adam",
				},
				body: &score.AddScore{
					Score: 100.2,
				},
			},
			fn: func(in args) *httpexpect.Object {
				var e error
				command := in.body.(*score.AddScore)
				command.ClientID = "adam"

				h.mockScoreUsecase.EXPECT().AddIgnoreDuplicate(gomock.Any(), command).Return(e).Times(1)

				return h.mockHTTP.POST("/api/v1/dup/score").
					WithHeaders(in.headers).
					WithJSON(in.body).
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status")
			},
			want: map[string]interface{}{
				"status": "ok",
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {

			expect := test.fn(test.args)
			for k, w := range test.want {
				expect.ValueEqual(k, w)
			}
		})
	}
}

// Test_GetLeaderBoard
func (h *handlerSuite) Test_GetLeaderBoard() {
	type args struct {
		headers map[string]string
	}

	tests := []struct {
		name string
		args args
		fn   func(args) *httpexpect.Object
		want map[string]interface{}
	}{
		{
			name: "test GetLeaderBoard occur error",
			args: args{},
			fn: func(args) *httpexpect.Object {
				h.mockScoreUsecase.EXPECT().GetLeaderBoard(gomock.Any()).Return(nil, errors.New("error")).Times(1)

				return h.mockHTTP.GET("/api/v1/leaderboard").
					Expect().
					Status(httptest.StatusOK).
					JSON().Object().
					ContainsKey("status").
					Value("status").Object()
			},
			want: map[string]interface{}{
				"message": "error",
			},
		},
		{
			name: "test GetLeaderBoard success",
			args: args{},
			fn: func(args) *httpexpect.Object {
				scores := []*model.Score{
					{
						ClientID: "adam",
						Score:    100.3,
					},
					{
						ClientID: "peter",
						Score:    10.1,
					},
				}
				h.mockScoreUsecase.EXPECT().GetLeaderBoard(gomock.Any()).Return(scores, nil).Times(1)

				return h.mockHTTP.GET("/api/v1/leaderboard").
					Expect().
					Status(httptest.StatusOK).
					JSON().Object()
			},
			want: map[string]interface{}{
				"topPlayers": []*model.Score{
					{
						ClientID: "adam",
						Score:    100.3,
					},
					{
						ClientID: "peter",
						Score:    10.1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {

			expect := test.fn(test.args)
			for k, w := range test.want {
				expect.ValueEqual(k, w)
			}
		})
	}
}
