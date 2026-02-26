package marker_types

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type markerTypeRepository struct {
	db *gorm.DB
}

func NewMarkerTypeRepository(db *gorm.DB) *markerTypeRepository {
	return &markerTypeRepository{
		db,
	}
}

var _ MarkerTypeRepository = &markerTypeRepository{}

func (r *markerTypeRepository) CreateMarkerType(ctx context.Context, markerType *models.MarkerType) (*models.MarkerType, error) {

	err := r.db.WithContext(ctx).Create(&markerType).Error
	if err != nil {
		return nil, err
	}

	return markerType, nil
}

func (r *markerTypeRepository) UpdateMarkerType(ctx context.Context, id int32, markerType *models.MarkerType) (*models.MarkerType, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&markerType).Updates(markerType).Error
	if err != nil {
		return nil, err
	}

	return markerType, nil
}

func (r *markerTypeRepository) UpdateMarkerTypeByUUID(ctx context.Context, uuid string, markerType *models.MarkerType) (*models.MarkerType, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&markerType).Updates(markerType).Error
	if err != nil {
		return nil, err
	}

	return markerType, nil
}

func (s *markerTypeRepository) DeleteMarkerType(ctx context.Context, id int32) error {

	markerType := &models.MarkerType{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(markerType).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *markerTypeRepository) DeleteMarkerTypeByUUID(ctx context.Context, uuid string) error {

	markerType := &models.MarkerType{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(markerType).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *markerTypeRepository) GetMarkerTypeByID(ctx context.Context, id int32) (*models.MarkerType, error) {

	var markerType = &models.MarkerType{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&markerType).Error
	if err != nil {
		return nil, err
	}

	return markerType, nil
}

func (r *markerTypeRepository) GetMarkerTypeByUUID(ctx context.Context, uuid string) (*models.MarkerType, error) {

	var markerType = &models.MarkerType{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&markerType).Error
	if err != nil {
		return nil, err
	}

	return markerType, nil
}

func (r *markerTypeRepository) GetAllMarkerTypes(ctx context.Context) ([]*models.MarkerType, error) {

	var markerTypes []*models.MarkerType

	err := r.db.WithContext(ctx).Find(&markerTypes).Error
	if err != nil {
		return nil, err
	}

	return markerTypes, nil

}

func (r *markerTypeRepository) PageMarkerType(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var markerTypes []models.MarkerType

	var err = r.db.Scopes(pkg.Paginate(markerTypes, &pagination, r.db)).Find(&markerTypes).Error
	pagination.Rows = make([]any, len(markerTypes))
	for i, markerType := range markerTypes {
		pagination.Rows[i] = markerType
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
