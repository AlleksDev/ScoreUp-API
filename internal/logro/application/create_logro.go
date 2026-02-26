package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/repository"
)

type CreateLogro struct {
	repo repository.LogroRepository
}

func NewCreateLogro(repo repository.LogroRepository) *CreateLogro {
	return &CreateLogro{repo: repo}
}

func (uc *CreateLogro) Execute(logro *entities.Logro) error {
	if logro.Name == "" {
		return fmt.Errorf("el nombre es obligatorio")
	}
	if logro.Description == "" {
		return fmt.Errorf("la descripción es obligatoria")
	}

	if err := uc.repo.Save(logro); err != nil {
		return fmt.Errorf("error al guardar logro: %v", err)
	}

	return nil
}
