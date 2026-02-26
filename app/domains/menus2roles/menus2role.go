package menus2roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type Menus2roleRepository interface {
	CreateMenus2role(ctx context.Context, menus2role *models.Menus2role) (*models.Menus2role, error)
	UpdateMenus2role(ctx context.Context, id int32, menus2role *models.Menus2role) (*models.Menus2role, error)
	UpdateMenus2roleByUUID(ctx context.Context, uuid string, menus2role *models.Menus2role) (*models.Menus2role, error)
	DeleteMenus2role(ctx context.Context, id int32) error
	DeleteMenus2roleByUUID(ctx context.Context, uuid string) error
	GetMenus2roleByID(ctx context.Context, id int32) (*models.Menus2role, error)
	GetMenus2roleByUUID(ctx context.Context, uuid string) (*models.Menus2role, error)
	GetMenus2roleByMenuUUID(ctx context.Context, menuUuid string) ([]*models.Menus2role, error)
	GetAllMenus2roles(ctx context.Context) ([]*models.Menus2role, error)
	PageMenus2role(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
