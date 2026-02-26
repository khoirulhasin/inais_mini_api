package drivers

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error)
	UpdateDriver(ctx context.Context, id int32, driver *models.Driver) (*models.Driver, error)
	UpdateDriverByUUID(ctx context.Context, uuid string, driver *models.Driver) (*models.Driver, error)
	DeleteDriver(ctx context.Context, id int32) error
	DeleteDriverByUUID(ctx context.Context, uuid string) error
	GetDriverByID(ctx context.Context, id int32) (*models.Driver, error)
	GetDriverByUUID(ctx context.Context, uuid string) (*models.Driver, error)
	GetAllDrivers(ctx context.Context) ([]*models.Driver, error)
	PageDriver(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
