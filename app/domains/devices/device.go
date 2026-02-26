package devices

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error)
	UpdateDevice(ctx context.Context, id int32, device *models.Device) (*models.Device, error)
	UpdateDeviceByUUID(ctx context.Context, uuid string, device *models.Device) (*models.Device, error)
	DeleteDevice(ctx context.Context, id int32) error
	DeleteDeviceByUUID(ctx context.Context, uuid string) error
	GetDeviceByID(ctx context.Context, id int32) (*models.Device, error)
	GetDeviceByUUID(ctx context.Context, uuid string) (*models.Device, error)
	GetAllDevices(ctx context.Context) ([]*models.Device, error)
	PageDevice(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
