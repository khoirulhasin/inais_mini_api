package menus2roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type menus2roleRepository struct {
	db *gorm.DB
}

func NewMenus2roleRepository(db *gorm.DB) *menus2roleRepository {
	return &menus2roleRepository{
		db,
	}
}

var _ Menus2roleRepository = &menus2roleRepository{}

func (r *menus2roleRepository) CreateMenus2role(ctx context.Context, menus2role *models.Menus2role) (*models.Menus2role, error) {

	err := r.db.WithContext(ctx).Create(&menus2role).Error
	if err != nil {
		return nil, err
	}

	return menus2role, nil
}

func (r *menus2roleRepository) UpdateMenus2role(ctx context.Context, id int32, menus2role *models.Menus2role) (*models.Menus2role, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&menus2role).Updates(menus2role).Error
	if err != nil {
		return nil, err
	}

	return menus2role, nil
}

func (r *menus2roleRepository) UpdateMenus2roleByUUID(ctx context.Context, uuid string, menus2role *models.Menus2role) (*models.Menus2role, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&menus2role).Updates(menus2role).Error
	if err != nil {
		return nil, err
	}

	return menus2role, nil
}

func (r *menus2roleRepository) DeleteMenus2role(ctx context.Context, id int32) error {

	menus2role := &models.Menus2role{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(menus2role).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *menus2roleRepository) DeleteMenus2roleByUUID(ctx context.Context, uuid string) error {

	menus2role := &models.Menus2role{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(menus2role).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *menus2roleRepository) GetMenus2roleByID(ctx context.Context, id int32) (*models.Menus2role, error) {

	var menus2role = &models.Menus2role{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&menus2role).Error
	if err != nil {
		return nil, err
	}

	return menus2role, nil
}

func (r *menus2roleRepository) GetMenus2roleByUUID(ctx context.Context, uuid string) (*models.Menus2role, error) {

	var menus2role = &models.Menus2role{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&menus2role).Error
	if err != nil {
		return nil, err
	}

	return menus2role, nil
}

func (r *menus2roleRepository) GetMenus2roleByMenuUUID(ctx context.Context, menuUuid string) ([]*models.Menus2role, error) {

	var menus2roles []*models.Menus2role

	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN menus ON menus2roles.menu_id = menus.id").
		Where("menus.uuid = ?", menuUuid).
		Preload("Menu").
		Preload("Role").
		Find(&menus2roles).Error

	if err != nil {
		return nil, err
	}

	return menus2roles, nil
}

func (r *menus2roleRepository) GetAllMenus2roles(ctx context.Context) ([]*models.Menus2role, error) {

	var menus2roles []*models.Menus2role

	err := r.db.WithContext(ctx).Find(&menus2roles).Error
	if err != nil {
		return nil, err
	}

	return menus2roles, nil

}

func (r *menus2roleRepository) PageMenus2role(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var menus2roles []models.Menus2role

	var err = r.db.WithContext(ctx).Scopes(pkg.Paginate(menus2roles, &pagination, r.db)).Find(&menus2roles).Error
	pagination.Rows = make([]any, len(menus2roles))
	for i, menus2role := range menus2roles {
		pagination.Rows[i] = menus2role
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
