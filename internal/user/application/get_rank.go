package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/user/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/user/domain/repository"
)

type GetRank struct {
	repo repository.UserRepository
}

func NewGetRank(repo repository.UserRepository) *GetRank {
	return &GetRank{repo: repo}
}

func (uc *GetRank) Execute() ([]entities.RankUser, error) {
	rank, err := uc.repo.GetRank()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ranking: %w", err)
	}
	return rank, nil
}