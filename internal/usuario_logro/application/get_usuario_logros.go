package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/repository"
)

type GetUsuarioLogros struct {
	repo repository.UsuarioLogroRepository
}

func NewGetUsuarioLogros(repo repository.UsuarioLogroRepository) *GetUsuarioLogros {
	return &GetUsuarioLogros{repo: repo}
}

func (uc *GetUsuarioLogros) Execute(userID int64) ([]*entities.UsuarioLogro, error) {
	logros, err := uc.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener logros del usuario: %v", err)
	}
	return logros, nil
}
