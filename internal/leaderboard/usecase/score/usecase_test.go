package score

import (
	"context"
	"errors"
	"leaderboard/internal/leaderboard/domain/model"
	"leaderboard/test/mock/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

// TestSuite
type TestSuite struct {
	suite.Suite
	ctrl                      *gomock.Controller
	mockLeaderBoardRepository *repository.MockLeaderBoardRepository
	usecase                   *usecase
}

// SetupTest
func (t *TestSuite) SetupSuite() {
	t.ctrl = gomock.NewController(t.T())
	t.mockLeaderBoardRepository = repository.NewMockLeaderBoardRepository(t.ctrl)

	t.usecase = &usecase{
		leaderBoardRepository: t.mockLeaderBoardRepository,
	}
}

// TestScoreUsecase
func TestScoreUsecase(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// Test_Add
func (t *TestSuite) Test_Add() {
	type args struct {
		ctx     context.Context
		command *AddScore
	}

	tests := []struct {
		name      string
		fn        func(args)
		args      args
		wantError bool
	}{
		{
			name: "test add score when leaderboard is empty case",
			fn: func(in args) {
				var result int64 = 0
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), "leaderboard").Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), "leaderboard", &model.Score{
					ClientID: in.command.ClientID,
					Score:    in.command.Score,
				}).Return(nil).Times(1)

				t.mockLeaderBoardRepository.EXPECT().SetExpire(gomock.Any(), "leaderboard", time.Minute*10).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "adam",
					Score:    10.2,
				},
			},
			wantError: false,
		},
		{
			name: "test add score when leaderboard not empty case",
			fn: func(in args) {
				var result int64 = 1
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), "leaderboard").Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), "leaderboard", &model.Score{
					ClientID: in.command.ClientID,
					Score:    in.command.Score,
				}).Return(nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "peter",
					Score:    91.2,
				},
			},
			wantError: false,
		},
		{
			name: "test add score error case",
			fn: func(in args) {
				var result int64 = 1
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), "leaderboard").Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), "leaderboard", &model.Score{
					ClientID: in.command.ClientID,
					Score:    in.command.Score,
				}).Return(errors.New("")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "Linda",
					Score:    91.2,
				},
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got := t.usecase.Add(test.args.ctx, test.args.command)
			t.Equal(test.wantError, got != nil)

		})
	}
}

// Test_AddIgnoreDuplicate
func (t *TestSuite) Test_AddIgnoreDuplicate() {
	type args struct {
		ctx     context.Context
		command *AddScore
	}

	tests := []struct {
		name      string
		fn        func(args)
		args      args
		wantError bool
	}{
		{
			name: "test add score when leaderboard is empty case",
			fn: func(in args) {
				var result int64 = 0
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), key).Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), key, gomock.Any()).Return(nil).Times(1)

				t.mockLeaderBoardRepository.EXPECT().SetExpire(gomock.Any(), key, time.Minute*10).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "adam",
					Score:    10.2,
				},
			},
			wantError: false,
		},
		{
			name: "test add score when leaderboard not empty case",
			fn: func(in args) {
				var result int64 = 1
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), key).Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), key, gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "peter",
					Score:    91.2,
				},
			},
			wantError: false,
		},
		{
			name: "test add score error case",
			fn: func(in args) {
				var result int64 = 1
				t.mockLeaderBoardRepository.EXPECT().Exists(gomock.Any(), key).Return(result).Times(1)

				t.mockLeaderBoardRepository.EXPECT().Create(gomock.Any(), key, gomock.Any()).Return(errors.New("")).Times(1)
			},
			args: args{
				ctx: context.Background(),
				command: &AddScore{
					ClientID: "John",
					Score:    91.2,
				},
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got := t.usecase.Add(test.args.ctx, test.args.command)
			t.Equal(test.wantError, got != nil)

		})
	}
}

// Test_GetLeaderBoard
func (t *TestSuite) Test_GetLeaderBoard() {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name       string
		fn         func(args)
		args       args
		wantResult []*model.Score
	}{
		{
			name: "test get leaderboard success case",
			fn: func(in args) {
				var (
					offset int64 = 0
					limit  int64 = 9
				)
				result := []*model.Score{
					{
						ClientID: "adam",
						Score:    100,
					},
					{
						ClientID: "peter",
						Score:    90,
					},
				}
				t.mockLeaderBoardRepository.EXPECT().List(gomock.Any(), key, offset, limit).Return(result, nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
			},
			wantResult: []*model.Score{
				{
					ClientID: "adam",
					Score:    100,
				},
				{
					ClientID: "peter",
					Score:    90,
				},
			},
		},
		{
			name: "test get leaderboard success case",
			fn: func(in args) {
				var (
					offset int64 = 0
					limit  int64 = 9
				)
				result := []*model.Score{
					{
						ClientID: `{"clientId":"linda","createdAt":1}`,
						Score:    100,
					},
					{
						ClientID: "peter",
						Score:    90,
					},
					{
						ClientID: `{"clientId":"adam","createdAt":2}`,
						Score:    8.2,
					},
				}
				t.mockLeaderBoardRepository.EXPECT().List(gomock.Any(), key, offset, limit).Return(result, nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
			},
			wantResult: []*model.Score{
				{
					ClientID: "linda",
					Score:    100,
				},
				{
					ClientID: "peter",
					Score:    90,
				},
				{
					ClientID: "adam",
					Score:    8.2,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got, _ := t.usecase.GetLeaderBoard(test.args.ctx)
			t.Equal(test.wantResult, got)
		})
	}
}

// Test_ResetLeaderBoard
func (t *TestSuite) Test_ResetLeaderBoard() {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name      string
		fn        func(args)
		args      args
		wantError bool
	}{
		{
			name: "test ResetLeaderBoard case",
			fn: func(in args) {
				t.mockLeaderBoardRepository.EXPECT().DeleteAll(in.ctx).Return(nil).Times(1)
			},
			args: args{
				ctx: context.Background(),
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got := t.usecase.ResetLeaderBoard(test.args.ctx)
			t.Equal(test.wantError, got != nil)

		})
	}
}
