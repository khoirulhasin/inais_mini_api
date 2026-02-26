package roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db,
	}
}

var _ RoleRepository = &roleRepository{}

func (r *roleRepository) CreateRole(ctx context.Context, role *models.Role) (*models.Role, error) {

	err := r.db.WithContext(ctx).Create(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, id int32, role *models.Role) (*models.Role, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&role).Updates(role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) UpdateRoleByUUID(ctx context.Context, uuid string, role *models.Role) (*models.Role, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&role).Updates(role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleRepository) DeleteRole(ctx context.Context, id int32) error {

	role := &models.Role{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(role).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *roleRepository) DeleteRoleByUUID(ctx context.Context, uuid string) error {

	role := &models.Role{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(role).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *roleRepository) GetRoleByID(ctx context.Context, id int32) (*models.Role, error) {

	var role = &models.Role{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) GetRoleByUUID(ctx context.Context, uuid string) (*models.Role, error) {

	var role = &models.Role{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) GetAllRoles(ctx context.Context) ([]*models.Role, error) {

	var roles []*models.Role

	err := r.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil

}

func (r *roleRepository) PageRole(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var roles []models.Role

	var err = r.db.Scopes(pkg.Paginate(roles, &pagination, r.db)).Find(&roles).Error
	pagination.Rows = make([]any, len(roles))
	for i, role := range roles {
		pagination.Rows[i] = role
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
