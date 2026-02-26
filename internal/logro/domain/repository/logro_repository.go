package repository

import "github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"

type LogroRepository interface {
	Save(logro *entities.Logro) error
	GetByID(id int64) (*entities.Logro, error)
	GetAll() ([]*entities.Logro, error)
	Update(logro *entities.Logro) error
	Delete(id int64) error
}
