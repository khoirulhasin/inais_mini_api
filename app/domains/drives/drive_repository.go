package drives

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type driveRepository struct {
	db *gorm.DB
}

func NewDriveRepository(db *gorm.DB) *driveRepository {
	return &driveRepository{
		db,
	}
}

var _ DriveRepository = &driveRepository{}

func (r *driveRepository) CreateDrive(ctx context.Context, drive *models.Drive) (*models.Drive, error) {

	err := r.db.WithContext(ctx).Create(&drive).Error
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (r *driveRepository) UpdateDrive(ctx context.Context, id int32, drive *models.Drive) (*models.Drive, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&drive).Updates(drive).Error
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (r *driveRepository) UpdateDriveByUUID(ctx context.Context, uuid string, drive *models.Drive) (*models.Drive, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&drive).Updates(drive).Error
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (s *driveRepository) DeleteDrive(ctx context.Context, id int32) error {

	drive := &models.Drive{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(drive).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *driveRepository) DeleteDriveByUUID(ctx context.Context, uuid string) error {

	drive := &models.Drive{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(drive).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *driveRepository) GetDriveByID(ctx context.Context, id int32) (*models.Drive, error) {

	var drive = &models.Drive{}

	err := r.db.WithContext(ctx).
		Preload("Ship").
		Preload("Driver").
		Where("id = ?", id).
		Take(&drive).Error
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (r *driveRepository) GetDriveByUUID(ctx context.Context, uuid string) (*models.Drive, error) {

	var drive = &models.Drive{}

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Preload("Ship").
		Preload("Driver").
		Take(&drive).Error
	if err != nil {
		return nil, err
	}

	return drive, nil
}

func (r *driveRepository) GetAllDrives(ctx context.Context) ([]*models.Drive, error) {

	var drives []*models.Drive

	err := r.db.WithContext(ctx).Find(&drives).Error
	if err != nil {
		return nil, err
	}

	return drives, nil

}

func (r *driveRepository) PageDrive(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var drives []models.Drive

	var err = r.db.WithContext(ctx).
		Preload("Ship").
		Preload("Driver").
		Scopes(pkg.Paginate(&models.Drive{}, &pagination, r.db)).
		Find(&drives).Error
	pagination.Rows = make([]any, len(drives))
	for i, drive := range drives {
		pagination.Rows[i] = drive
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
