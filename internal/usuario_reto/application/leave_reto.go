package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/repository"
)

type LeaveReto struct {
	repo repository.UsuarioRetoRepository
}

func NewLeaveReto(repo repository.UsuarioRetoRepository) *LeaveReto {
	return &LeaveReto{repo: repo}
}

func (uc *LeaveReto) Execute(userID int64, retoID int64) error {
	// Verificar que esté unido
	existing, err := uc.repo.GetByUserAndReto(userID, retoID)
	if err != nil {
		return fmt.Errorf("error al verificar unión: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("el usuario no está unido a este reto")
	}

	if err := uc.repo.Delete(userID, retoID); err != nil {
		return fmt.Errorf("error al abandonar reto: %v", err)
	}

	return nil
}
