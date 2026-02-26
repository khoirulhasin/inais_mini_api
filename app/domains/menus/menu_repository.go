package menus

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *menuRepository {
	return &menuRepository{
		db,
	}
}

var _ MenuRepository = &menuRepository{}

func (r *menuRepository) CreateMenu(ctx context.Context, menu *models.Menu) (*models.Menu, error) {

	err := r.db.WithContext(ctx).Create(&menu).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *menuRepository) UpdateMenu(ctx context.Context, id int32, menu *models.Menu) (*models.Menu, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&menu).Updates(menu).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *menuRepository) UpdateMenuByUUID(ctx context.Context, uuid string, menu *models.Menu) (*models.Menu, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&menu).Updates(menu).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *menuRepository) DeleteMenu(ctx context.Context, id int32) error {

	menu := &models.Menu{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(menu).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *menuRepository) DeleteMenuByUUID(ctx context.Context, uuid string) error {

	menu := &models.Menu{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(menu).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *menuRepository) GetMenuByID(ctx context.Context, id int32) (*models.Menu, error) {

	var menu = &models.Menu{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&menu).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *menuRepository) GetMenuByUUID(ctx context.Context, uuid string) (*models.Menu, error) {

	var menu = &models.Menu{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&menu).Error
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *menuRepository) GetAllMenus(ctx context.Context) ([]*models.Menu, error) {

	var menus []*models.Menu

	err := r.db.WithContext(ctx).Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil

}

func (r *menuRepository) GetMenuParent(ctx context.Context, roleId int) ([]*models.Menu, error) {
	var menus []*models.Menu

	var err = r.db.
		WithContext(ctx).
		Joins("LEFT JOIN menus2roles on menus.id = menus2roles.menu_id").
		Where("level = ?", 0).
		Where("show = ?", 1).
		Where("menus2roles.role_id =  ?", roleId).
		Order("sequence ASC").
		Preload("Menus", preload).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) GetMenuFlat(ctx context.Context, roleId int) ([]*models.Menu, error) {
	var menus []*models.Menu

	var err = r.db.
		WithContext(ctx).
		Joins("LEFT JOIN menus2roles on menus.id = menus2roles.menu_id").
		Where("menus2roles.role_id =  ?", roleId).
		Order("sequence asc").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (u menuRepository) GetMenuAllParents(ctx context.Context) ([]*models.Menu, error) {
	var menus []*models.Menu

	var err = u.db.
		WithContext(ctx).
		Where("level = ?", 0).
		Order("sequence asc").
		Preload(clause.Associations, preload).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) PageMenu(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var menus []models.Menu

	var err = r.db.WithContext(ctx).Scopes(pkg.Paginate(menus, &pagination, r.db)).Find(&menus).Error
	pagination.Rows = make([]any, len(menus))
	for i, menu := range menus {
		pagination.Rows[i] = menu
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
