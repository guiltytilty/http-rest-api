package store

import (
	"context"
	"http-rest-api/internal/app/model"
)

type UserRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*model.User, error)
	Create(ctx context.Context, u *model.User) error
	FindByID(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, upd model.UserUpdate) (*model.User, error)
	Delete(ctx context.Context, id int) error
}
