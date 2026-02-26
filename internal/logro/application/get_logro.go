package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/repository"
)

type GetLogro struct {
	repo repository.LogroRepository
}

func NewGetLogro(repo repository.LogroRepository) *GetLogro {
	return &GetLogro{repo: repo}
}

func (uc *GetLogro) ExecuteByID(id int64) (*entities.Logro, error) {
	logro, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener logro: %v", err)
	}
	if logro == nil {
		return nil, fmt.Errorf("logro no encontrado")
	}
	return logro, nil
}

func (uc *GetLogro) ExecuteAll() ([]*entities.Logro, error) {
	logros, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener logros: %v", err)
	}
	return logros, nil
}
