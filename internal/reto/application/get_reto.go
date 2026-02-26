package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/repository"
)

type GetReto struct {
	repo repository.RetoRepository
}

func NewGetReto(repo repository.RetoRepository) *GetReto {
	return &GetReto{repo: repo}
}

func (uc *GetReto) ExecuteByID(id int64) (*entities.Reto, error) {
	reto, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reto: %v", err)
	}
	if reto == nil {
		return nil, fmt.Errorf("reto no encontrado")
	}
	return reto, nil
}

func (uc *GetReto) ExecuteAll() ([]*entities.Reto, error) {
	retos, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener retos: %v", err)
	}
	return retos, nil
}

func (uc *GetReto) ExecuteByCreator(userID int64) ([]*entities.Reto, error) {
	retos, err := uc.repo.GetByCreator(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener retos del creador: %v", err)
	}
	return retos, nil
}
