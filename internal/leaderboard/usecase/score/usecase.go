package score

import (
	"context"
	"leaderboard/internal/leaderboard/domain/model"
	"leaderboard/internal/leaderboard/domain/repository"
	"leaderboard/pkg/encoder/json"
	"time"
)

const (
	key = "leaderboard"
)

type usecase struct {
	leaderBoardRepository repository.LeaderBoardRepository
}

// NewUseCase -
func NewUseCase(leaderBoardRepository repository.LeaderBoardRepository) ScoreUsecase {
	return &usecase{
		leaderBoardRepository: leaderBoardRepository,
	}
}

// Add - add one score record
func (u *usecase) Add(ctx context.Context, command *AddScore) error {
	in := &model.Score{
		ClientID: command.ClientID,
		Score:    command.Score,
	}

	// check if key exists
	setExpire := false
	if u.leaderBoardRepository.Exists(ctx, key) == 0 {
		setExpire = true
	}

	err := u.leaderBoardRepository.Create(ctx, key, in)
	if err != nil {
		return err
	}

	// If it is set TTL(10 minute)
	if setExpire {
		u.leaderBoardRepository.SetExpire(ctx, key, time.Minute*10)
	}

	return nil
}

// AddIgnoreDuplicate - the duplicate ClientID can appear on leaderboard
func (u *usecase) AddIgnoreDuplicate(ctx context.Context, command *AddScore) error {
	s := &model.Score{
		ClientID:  command.ClientID,
		CreatedAt: time.Now().Unix(),
	}

	coder := json.NewEncoder()
	c, _ := coder.Encode(&s)

	in := &model.Score{
		ClientID: string(c),
		Score:    command.Score,
	}

	// check if key exists
	setExpire := false
	if u.leaderBoardRepository.Exists(ctx, key) == 0 {
		setExpire = true
	}

	err := u.leaderBoardRepository.Create(ctx, key, in)
	if err != nil {
		return err
	}

	// If it is set TTL(10 minute)
	if setExpire {
		u.leaderBoardRepository.SetExpire(ctx, key, time.Minute*10)
	}

	return nil
}

// GetLeaderBoard
func (u *usecase) GetLeaderBoard(ctx context.Context) ([]*model.Score, error) {
	// get leaderboard for top 10(0 - 9) with score
	scores, err := u.leaderBoardRepository.List(ctx, key, 0, 9)
	if err != nil {
		return nil, err
	}

	coder := json.NewEncoder()
	for i, v := range scores {
		s := &model.Score{}
		if err := coder.Decode([]byte(v.ClientID), &s); err != nil {
			continue
		}
		scores[i].ClientID = s.ClientID
	}

	return scores, nil
}

// ResetLeaderBoard
func (u *usecase) ResetLeaderBoard(ctx context.Context) error {
	return u.leaderBoardRepository.DeleteAll(ctx)
}
