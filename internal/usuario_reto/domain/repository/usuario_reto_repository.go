package repository

import "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/entities"

type UsuarioRetoRepository interface {
	Save(ur *entities.UsuarioReto) error
	GetByUserID(userID int64) ([]*entities.UsuarioReto, error)
	GetByRetoID(retoID int64) ([]*entities.UsuarioReto, error)
	GetByUserAndReto(userID int64, retoID int64) (*entities.UsuarioReto, error)
	Update(ur *entities.UsuarioReto) error
	Delete(userID int64, retoID int64) error
}
