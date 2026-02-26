package ships

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type shipRepository struct {
	db *gorm.DB
}

func NewShipRepository(db *gorm.DB) *shipRepository {
	return &shipRepository{
		db,
	}
}

var _ ShipRepository = &shipRepository{}

func (r *shipRepository) CreateShip(ctx context.Context, Ship *models.Ship) (*models.Ship, error) {

	err := r.db.WithContext(ctx).Create(&Ship).Error
	if err != nil {
		return nil, err
	}

	return Ship, nil
}

func (r *shipRepository) UpdateShip(ctx context.Context, id int32, Ship *models.Ship) (*models.Ship, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&Ship).Updates(Ship).Error
	if err != nil {
		return nil, err
	}

	return Ship, nil
}

func (r *shipRepository) UpdateShipByUUID(ctx context.Context, uuid string, Ship *models.Ship) (*models.Ship, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&Ship).Updates(Ship).Error
	if err != nil {
		return nil, err
	}

	return Ship, nil
}

func (r *shipRepository) DeleteShip(ctx context.Context, id int32) error {

	Ship := &models.Ship{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(Ship).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *shipRepository) DeleteShipByUUID(ctx context.Context, uuid string) error {

	Ship := &models.Ship{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(Ship).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *shipRepository) GetShipByID(ctx context.Context, id int32) (*models.Ship, error) {

	var Ship = &models.Ship{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&Ship).Error
	if err != nil {
		return nil, err
	}

	return Ship, nil
}

func (r *shipRepository) GetShipByUUID(ctx context.Context, uuid string) (*models.Ship, error) {

	var Ship = &models.Ship{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&Ship).Error
	if err != nil {
		return nil, err
	}

	return Ship, nil
}

func (r *shipRepository) GetAllShips(ctx context.Context) ([]*models.Ship, error) {

	var Ships []*models.Ship

	err := r.db.WithContext(ctx).
		Find(&Ships).Error
	if err != nil {
		return nil, err
	}

	return Ships, nil

}

func (r *shipRepository) PageShip(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var Ships []models.Ship

	var err = r.db.WithContext(ctx).
		Scopes(pkg.Paginate(models.Ship{}, &pagination, r.db)).
		Find(&Ships).Error
	pagination.Rows = make([]any, len(Ships))
	for i, Ship := range Ships {

		pagination.Rows[i] = Ship
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
