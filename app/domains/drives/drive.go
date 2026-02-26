package drives

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type DriveRepository interface {
	CreateDrive(ctx context.Context, drive *models.Drive) (*models.Drive, error)
	UpdateDrive(ctx context.Context, id int32, drive *models.Drive) (*models.Drive, error)
	UpdateDriveByUUID(ctx context.Context, uuid string, drive *models.Drive) (*models.Drive, error)
	DeleteDrive(ctx context.Context, id int32) error
	DeleteDriveByUUID(ctx context.Context, uuid string) error
	GetDriveByID(ctx context.Context, id int32) (*models.Drive, error)
	GetDriveByUUID(ctx context.Context, uuid string) (*models.Drive, error)
	GetAllDrives(ctx context.Context) ([]*models.Drive, error)
	PageDrive(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
