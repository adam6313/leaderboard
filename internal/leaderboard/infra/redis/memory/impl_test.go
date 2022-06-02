package memory

import (
	"context"
	"errors"
	"leaderboard/internal/leaderboard/domain/model"
	"testing"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
)

// TestSuite
type TestSuite struct {
	suite.Suite
	mockClient redismock.ClientMock
	Repo       *Repo
}

// SetupTest
func (t *TestSuite) SetupSuite() {
	client, mockClient := redismock.NewClientMock()
	t.mockClient = mockClient

	t.Repo = &Repo{
		client: client,
	}
}

// TestLeaderBoardRepository
func TestLeaderBoardRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// Create
func (t *TestSuite) Test_Create() {
	type args struct {
		ctx   context.Context
		key   string
		score *model.Score
	}

	tests := []struct {
		name      string
		fn        func(args)
		args      args
		wantError bool
	}{
		{
			name: "test zadd success",
			fn: func(in args) {
				z := &goredis.Z{
					Member: in.score.ClientID,
					Score:  in.score.Score,
				}
				t.mockClient.ExpectZAdd(in.key, z).SetVal(1)
			},
			args: args{
				ctx: context.Background(),
				key: "leaderboard",
				score: &model.Score{
					ClientID: "test_adam",
					Score:    5.1,
				},
			},
			wantError: false,
		},
		{
			name: "test member is empty",
			fn: func(in args) {
				z := &goredis.Z{
					Score: in.score.Score,
				}
				t.mockClient.ExpectZAdd(in.key, z)
			},
			args: args{
				ctx: context.Background(),
				key: "leaderboard",
				score: &model.Score{
					Score: 5.1,
				},
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got := t.Repo.Create(test.args.ctx, test.args.key, test.args.score)
			t.Equal(test.wantError, got != nil)
			t.mockClient.ClearExpect()
		})
	}
}

// Test_List
func (t *TestSuite) Test_List() {
	type args struct {
		ctx    context.Context
		key    string
		offset int64
		limit  int64
	}

	tests := []struct {
		name       string
		fn         func(args)
		args       args
		wantResult []*model.Score
		wantError  bool
	}{
		{
			name: "test get list case",
			fn: func(in args) {
				res := []goredis.Z{
					{
						Member: "a",
						Score:  40,
					},
					{
						Member: "b",
						Score:  30,
					},
				}

				t.mockClient.ExpectZRevRangeWithScores(in.key, in.offset, in.limit).SetVal(res)
			},
			args: args{
				ctx:    context.Background(),
				key:    "leaderboard",
				offset: 0,
				limit:  9,
			},
			wantResult: []*model.Score{
				{
					ClientID: "a",
					Score:    40,
				},
				{
					ClientID: "b",
					Score:    30,
				},
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got, err := t.Repo.List(test.args.ctx, test.args.key, test.args.offset, test.args.limit)
			t.Equal(test.wantError, err != nil)
			t.Equal(got, test.wantResult)

			t.mockClient.ClearExpect()
		})
	}
}

// Test_DeleteAll
func (t *TestSuite) Test_DeleteAll() {
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
			name: "test deleteAll case",
			fn: func(in args) {
				t.mockClient.ExpectScan(0, "*", 0)
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

			err := t.Repo.DeleteAll(test.args.ctx)
			t.Equal(test.wantError, err != nil)

			t.mockClient.ClearExpect()
		})
	}
}

// Test_SetExpire
func (t *TestSuite) Test_SetExpire() {
	type args struct {
		ctx  context.Context
		key  string
		time time.Duration
	}

	tests := []struct {
		name      string
		fn        func(args)
		args      args
		wantError bool
	}{
		{
			name: "test SetExpire success case",
			fn: func(in args) {
				t.mockClient.ExpectExpire(in.key, in.time).SetVal(true)
			},
			args: args{
				ctx:  context.Background(),
				key:  "leaderboard",
				time: time.Minute * 10,
			},
			wantError: false,
		},
		{
			name: "test SetExpire error case",
			fn: func(in args) {
				t.mockClient.ExpectExpire(in.key, in.time).SetErr(errors.New(""))
			},
			args: args{
				ctx:  context.Background(),
				key:  "leaderboard",
				time: time.Minute * 10,
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			err := t.Repo.SetExpire(test.args.ctx, test.args.key, test.args.time)
			t.Equal(test.wantError, err != nil)

			t.mockClient.ClearExpect()
		})
	}
}

// Test_Exists
func (t *TestSuite) Test_Exists() {
	type args struct {
		ctx context.Context
		key string
	}

	tests := []struct {
		name       string
		fn         func(args)
		args       args
		wantResult int64
	}{
		{
			name: "test Exists success case",
			fn: func(in args) {
				t.mockClient.ExpectExists(in.key).SetVal(1)
			},
			args: args{
				ctx: context.Background(),
				key: "leaderboard",
			},
			wantResult: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.fn(test.args)

			got := t.Repo.Exists(test.args.ctx, test.args.key)
			t.Equal(test.wantResult, got)

			t.mockClient.ClearExpect()
		})
	}
}
