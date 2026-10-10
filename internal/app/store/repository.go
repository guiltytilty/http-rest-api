package store

import "http-rest-api/internal/app/model"

type UserRepository interface {
	GetAll(limit, offset int) ([]*model.User, error)
	Create(*model.User) error
	FindByID(int) (*model.User, error)
	Update(id int, upd model.UserUpdate) (*model.User, error)
	Delete(id int) error
}
