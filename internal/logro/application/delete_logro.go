package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/repository"
)

type DeleteLogro struct {
	repo repository.LogroRepository
}

func NewDeleteLogro(repo repository.LogroRepository) *DeleteLogro {
	return &DeleteLogro{repo: repo}
}

func (uc *DeleteLogro) Execute(id int64) error {
	existing, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error al buscar logro: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("logro no encontrado")
	}

	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar logro: %v", err)
	}

	return nil
}
