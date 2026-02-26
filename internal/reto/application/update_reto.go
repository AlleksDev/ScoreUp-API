package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/repository"
)

type UpdateReto struct {
	repo repository.RetoRepository
}

func NewUpdateReto(repo repository.RetoRepository) *UpdateReto {
	return &UpdateReto{repo: repo}
}

func (uc *UpdateReto) Execute(reto *entities.Reto) error {
	existing, err := uc.repo.GetByID(reto.ID)
	if err != nil {
		return fmt.Errorf("error al buscar reto: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("reto no encontrado")
	}

	if err := uc.repo.Update(reto); err != nil {
		return fmt.Errorf("error al actualizar reto: %v", err)
	}

	return nil
}
