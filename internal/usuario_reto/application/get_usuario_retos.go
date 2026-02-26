package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/repository"
)

type GetUsuarioRetos struct {
	repo repository.UsuarioRetoRepository
}

func NewGetUsuarioRetos(repo repository.UsuarioRetoRepository) *GetUsuarioRetos {
	return &GetUsuarioRetos{repo: repo}
}

func (uc *GetUsuarioRetos) ExecuteByUser(userID int64) ([]*entities.UsuarioReto, error) {
	results, err := uc.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener retos del usuario: %v", err)
	}
	return results, nil
}

func (uc *GetUsuarioRetos) ExecuteByReto(retoID int64) ([]*entities.UsuarioReto, error) {
	results, err := uc.repo.GetByRetoID(retoID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios del reto: %v", err)
	}
	return results, nil
}
