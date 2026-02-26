package drivers

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) *driverRepository {
	return &driverRepository{
		db,
	}
}

var _ DriverRepository = &driverRepository{}

func (r *driverRepository) CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {

	err := r.db.WithContext(ctx).Create(&driver).Error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (r *driverRepository) UpdateDriver(ctx context.Context, id int32, driver *models.Driver) (*models.Driver, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&driver).Updates(driver).Error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (r *driverRepository) UpdateDriverByUUID(ctx context.Context, uuid string, driver *models.Driver) (*models.Driver, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&driver).Updates(driver).Error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (r *driverRepository) DeleteDriver(ctx context.Context, id int32) error {

	driver := &models.Driver{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(driver).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *driverRepository) DeleteDriverByUUID(ctx context.Context, uuid string) error {

	driver := &models.Driver{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(driver).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *driverRepository) GetDriverByID(ctx context.Context, id int32) (*models.Driver, error) {

	var driver = &models.Driver{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&driver).Error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (r *driverRepository) GetDriverByUUID(ctx context.Context, uuid string) (*models.Driver, error) {

	var driver = &models.Driver{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&driver).Error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (r *driverRepository) GetAllDrivers(ctx context.Context) ([]*models.Driver, error) {

	var drivers []*models.Driver

	err := r.db.WithContext(ctx).Find(&drivers).Error
	if err != nil {
		return nil, err
	}

	return drivers, nil

}

func (r *driverRepository) PageDriver(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var drivers []models.Driver

	var err = r.db.WithContext(ctx).
		Scopes(pkg.Paginate(models.Driver{}, &pagination, r.db)).
		Find(&drivers).Error
	pagination.Rows = make([]any, len(drivers))
	for i, driver := range drivers {
		pagination.Rows[i] = driver
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
