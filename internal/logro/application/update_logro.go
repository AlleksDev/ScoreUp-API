package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/repository"
)

type UpdateLogro struct {
	repo repository.LogroRepository
}

func NewUpdateLogro(repo repository.LogroRepository) *UpdateLogro {
	return &UpdateLogro{repo: repo}
}

func (uc *UpdateLogro) Execute(logro *entities.Logro) error {
	existing, err := uc.repo.GetByID(logro.ID)
	if err != nil {
		return fmt.Errorf("error al buscar logro: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("logro no encontrado")
	}

	if err := uc.repo.Update(logro); err != nil {
		return fmt.Errorf("error al actualizar logro: %v", err)
	}

	return nil
}
