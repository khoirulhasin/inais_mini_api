package devices

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *deviceRepository {
	return &deviceRepository{
		db,
	}
}

var _ DeviceRepository = &deviceRepository{}

func (r *deviceRepository) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {

	err := r.db.WithContext(ctx).Create(&device).Error
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, id int32, device *models.Device) (*models.Device, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&device).Updates(device).Error
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) UpdateDeviceByUUID(ctx context.Context, uuid string, device *models.Device) (*models.Device, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&device).Updates(device).Error
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *deviceRepository) DeleteDevice(ctx context.Context, id int32) error {

	device := &models.Device{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(device).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *deviceRepository) DeleteDeviceByUUID(ctx context.Context, uuid string) error {

	device := &models.Device{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(device).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, id int32) (*models.Device, error) {

	var device = &models.Device{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&device).Error
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) GetDeviceByUUID(ctx context.Context, uuid string) (*models.Device, error) {

	var device = &models.Device{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&device).Error
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) GetAllDevices(ctx context.Context) ([]*models.Device, error) {

	var devices []*models.Device

	err := r.db.WithContext(ctx).Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil

}

func (r *deviceRepository) PageDevice(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var devices []models.Device

	var err = r.db.
		WithContext(ctx).
		Preload("Owner").
		Scopes(pkg.Paginate(&models.Device{}, &pagination, r.db)).
		Find(&devices).Error
	pagination.Rows = make([]any, len(devices))
	for i, device := range devices {
		pagination.Rows[i] = device
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
