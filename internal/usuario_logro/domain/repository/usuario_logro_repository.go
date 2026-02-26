package repository

import "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/entities"

type UsuarioLogroRepository interface {
	Save(ul *entities.UsuarioLogro) error
	GetByUserID(userID int64) ([]*entities.UsuarioLogro, error)
	Exists(userID int64, logroID int64) (bool, error)
	Delete(userID int64, logroID int64) error
}
