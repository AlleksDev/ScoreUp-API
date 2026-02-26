package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/repository"
)

type JoinReto struct {
	repo repository.UsuarioRetoRepository
}

func NewJoinReto(repo repository.UsuarioRetoRepository) *JoinReto {
	return &JoinReto{repo: repo}
}

func (uc *JoinReto) Execute(userID int64, retoID int64) error {
	// Verificar si ya está unido
	existing, err := uc.repo.GetByUserAndReto(userID, retoID)
	if err != nil {
		return fmt.Errorf("error al verificar unión existente: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("el usuario ya está unido a este reto")
	}

	ur := &entities.UsuarioReto{
		UserID:   userID,
		RetoID:   retoID,
		Progress: 0,
		Status:   "activo",
	}

	if err := uc.repo.Save(ur); err != nil {
		return fmt.Errorf("error al unirse al reto: %v", err)
	}

	return nil
}
