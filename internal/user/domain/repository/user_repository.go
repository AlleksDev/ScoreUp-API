package repository

import "github.com/AlleksDev/ScoreUp-API/internal/user/domain/entities"

type UserRepository interface {
	Save(user *entities.User) error
	GetByID(id int64) (*entities.User, error)
	GetUsersByIDs(ids []int64) ([]entities.User, error)
	Update(user *entities.User) error
	Delete(id int64) error
}