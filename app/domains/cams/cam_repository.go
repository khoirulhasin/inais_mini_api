package cams

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type camRepository struct {
	db *gorm.DB
}

func NewCamRepository(db *gorm.DB) *camRepository {
	return &camRepository{
		db,
	}
}

var _ CamRepository = &camRepository{}

func (r *camRepository) CreateCam(ctx context.Context, cam *models.Cam) (*models.Cam, error) {

	err := r.db.WithContext(ctx).Create(&cam).Error
	if err != nil {
		return nil, err
	}

	return cam, nil
}

func (r *camRepository) UpdateCam(ctx context.Context, id int32, cam *models.Cam) (*models.Cam, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&cam).Updates(cam).Error
	if err != nil {
		return nil, err
	}

	return cam, nil
}

func (r *camRepository) UpdateCamByUUID(ctx context.Context, uuid string, cam *models.Cam) (*models.Cam, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&cam).Updates(cam).Error
	if err != nil {
		return nil, err
	}

	return cam, nil
}

func (s *camRepository) DeleteCam(ctx context.Context, id int32) error {

	cam := &models.Cam{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(cam).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *camRepository) DeleteCamByUUID(ctx context.Context, uuid string) error {

	cam := &models.Cam{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(cam).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *camRepository) GetCamByID(ctx context.Context, id int32) (*models.Cam, error) {

	var cam = &models.Cam{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&cam).Error
	if err != nil {
		return nil, err
	}

	return cam, nil
}

func (r *camRepository) GetCamByUUID(ctx context.Context, uuid string) (*models.Cam, error) {

	var cam = &models.Cam{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&cam).Error
	if err != nil {
		return nil, err
	}

	return cam, nil
}

func (r *camRepository) GetCamByStateID(ctx context.Context, stateId int32) ([]*models.Cam, error) {

	var cams []*models.Cam

	err := r.db.WithContext(ctx).Where("state_id", stateId).Find(&cams).Error
	if err != nil {
		return nil, err
	}

	return cams, nil

}

func (r *camRepository) GetAllCams(ctx context.Context) ([]*models.Cam, error) {

	var cams []*models.Cam

	err := r.db.WithContext(ctx).Find(&cams).Error
	if err != nil {
		return nil, err
	}

	return cams, nil

}

func (r *camRepository) PageCam(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var cams []models.Cam

	var err = r.db.WithContext(ctx).
		Scopes(pkg.Paginate(&models.Cam{}, &pagination, r.db)).
		Find(&cams).Error
	pagination.Rows = make([]any, len(cams))
	for i, cam := range cams {
		pagination.Rows[i] = cam
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
