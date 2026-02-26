package menus

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type MenuRepository interface {
	CreateMenu(ctx context.Context, country *models.Menu) (*models.Menu, error)
	UpdateMenu(ctx context.Context, id int32, country *models.Menu) (*models.Menu, error)
	UpdateMenuByUUID(ctx context.Context, uuid string, country *models.Menu) (*models.Menu, error)
	DeleteMenu(ctx context.Context, id int32) error
	DeleteMenuByUUID(ctx context.Context, uuid string) error
	GetMenuByID(ctx context.Context, id int32) (*models.Menu, error)
	GetMenuByUUID(ctx context.Context, uuid string) (*models.Menu, error)
	GetAllMenus(ctx context.Context) ([]*models.Menu, error)
	GetMenuParent(ctx context.Context, roleId int) ([]*models.Menu, error)
	GetMenuFlat(ctx context.Context, roleId int) ([]*models.Menu, error)
	GetMenuAllParents(ctx context.Context) ([]*models.Menu, error)
	PageMenu(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}

func preload(d *gorm.DB) *gorm.DB {
	return d.Preload("Menus", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC").Preload("Menus", preload) // Recursive preload with ordering
	})
}
