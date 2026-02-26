package repository

import "github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"

type RetoRepository interface {
	Save(reto *entities.Reto) error
	GetByID(id int64) (*entities.Reto, error)
	GetAll() ([]*entities.Reto, error)
	GetByCreator(userID int64) ([]*entities.Reto, error)
	Update(reto *entities.Reto) error
	Delete(id int64) error
}
