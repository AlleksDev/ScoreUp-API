package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/repository"
)

type DeleteUsuarioLogro struct {
	repo repository.UsuarioLogroRepository
}

func NewDeleteUsuarioLogro(repo repository.UsuarioLogroRepository) *DeleteUsuarioLogro {
	return &DeleteUsuarioLogro{repo: repo}
}

func (uc *DeleteUsuarioLogro) Execute(userID int64, logroID int64) error {
	exists, err := uc.repo.Exists(userID, logroID)
	if err != nil {
		return fmt.Errorf("error al verificar logro del usuario: %v", err)
	}
	if !exists {
		return fmt.Errorf("el usuario no tiene este logro asignado")
	}

	if err := uc.repo.Delete(userID, logroID); err != nil {
		return fmt.Errorf("error al eliminar logro del usuario: %v", err)
	}

	return nil
}
