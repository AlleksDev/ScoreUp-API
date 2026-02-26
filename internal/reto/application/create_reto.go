package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/repository"
)

type CreateReto struct {
	repo repository.RetoRepository
}

func NewCreateReto(repo repository.RetoRepository) *CreateReto {
	return &CreateReto{repo: repo}
}

func (uc *CreateReto) Execute(reto *entities.Reto) error {
	if reto.Subject == "" {
		return fmt.Errorf("la materia es obligatoria")
	}
	if reto.Description == "" {
		return fmt.Errorf("la descripción es obligatoria")
	}
	if reto.Goal <= 0 {
		return fmt.Errorf("la meta debe ser mayor a 0")
	}

	if err := uc.repo.Save(reto); err != nil {
		return fmt.Errorf("error al guardar reto: %v", err)
	}

	return nil
}
