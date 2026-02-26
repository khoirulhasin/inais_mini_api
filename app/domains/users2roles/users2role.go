package users2roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type Users2roleRepository interface {
	CreateUsers2role(ctx context.Context, users2role *models.Users2role) (*models.Users2role, error)
	UpdateUsers2role(ctx context.Context, id int32, users2role *models.Users2role) (*models.Users2role, error)
	UpdateUsers2roleByUUID(ctx context.Context, uuid string, users2role *models.Users2role) (*models.Users2role, error)
	DeleteUsers2role(ctx context.Context, id int32) error
	DeleteUsers2roleByUUID(ctx context.Context, uuid string) error
	GetUsers2roleByID(ctx context.Context, id int32) (*models.Users2role, error)
	GetUsers2roleByUUID(ctx context.Context, uuid string) (*models.Users2role, error)
	GetUsers2roleByUserUUID(ctx context.Context, userUuid string) ([]*models.Users2role, error)
	GetUsers2roleByRoleId(ctx context.Context, roleId int32) ([]*models.Users2role, error)
	GetAllUsers2roles(ctx context.Context) ([]*models.Users2role, error)
	PageUsers2role(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
