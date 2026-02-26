package repository

import "github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"

type RetoRepository interface {
	Save(reto *entities.Reto) error
	GetByID(id int64) (*entities.Reto, error)
	GetByUserID(userID int64) ([]*entities.Reto, error)
	Update(reto *entities.Reto) error
	Delete(id int64) error
}
