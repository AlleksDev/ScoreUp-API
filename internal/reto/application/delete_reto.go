package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/repository"
)

type DeleteReto struct {
	repo repository.RetoRepository
}

func NewDeleteReto(repo repository.RetoRepository) *DeleteReto {
	return &DeleteReto{repo: repo}
}

func (uc *DeleteReto) Execute(id int64) error {
	existing, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error al buscar reto: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("reto no encontrado")
	}

	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar reto: %v", err)
	}

	return nil
}
