package users2roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type users2roleRepository struct {
	db *gorm.DB
}

func NewUsers2roleRepository(db *gorm.DB) *users2roleRepository {
	return &users2roleRepository{
		db,
	}
}

var _ Users2roleRepository = &users2roleRepository{}

func (r *users2roleRepository) CreateUsers2role(ctx context.Context, users2role *models.Users2role) (*models.Users2role, error) {

	err := r.db.WithContext(ctx).Create(&users2role).Error
	if err != nil {
		return nil, err
	}

	return users2role, nil
}

func (r *users2roleRepository) UpdateUsers2role(ctx context.Context, id int32, users2role *models.Users2role) (*models.Users2role, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&users2role).Updates(users2role).Error
	if err != nil {
		return nil, err
	}

	return users2role, nil
}

func (r *users2roleRepository) UpdateUsers2roleByUUID(ctx context.Context, uuid string, users2role *models.Users2role) (*models.Users2role, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&users2role).Updates(users2role).Error
	if err != nil {
		return nil, err
	}

	return users2role, nil
}

func (s *users2roleRepository) DeleteUsers2role(ctx context.Context, id int32) error {

	users2role := &models.Users2role{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(users2role).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *users2roleRepository) DeleteUsers2roleByUUID(ctx context.Context, uuid string) error {

	users2role := &models.Users2role{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(users2role).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *users2roleRepository) GetUsers2roleByID(ctx context.Context, id int32) (*models.Users2role, error) {

	var users2role = &models.Users2role{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&users2role).Error
	if err != nil {
		return nil, err
	}

	return users2role, nil
}

func (r *users2roleRepository) GetUsers2roleByUUID(ctx context.Context, uuid string) (*models.Users2role, error) {

	var users2role = &models.Users2role{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&users2role).Error
	if err != nil {
		return nil, err
	}

	return users2role, nil
}

func (r *users2roleRepository) GetUsers2roleByUserUUID(ctx context.Context, userUuid string) ([]*models.Users2role, error) {

	var users2roles []*models.Users2role

	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN users ON users2roles.user_id = users.id").
		Where("users.uuid = ?", userUuid).
		Preload("User").
		Preload("Role").
		Find(&users2roles).Error

	if err != nil {
		return nil, err
	}

	return users2roles, nil
}

func (r *users2roleRepository) GetUsers2roleByRoleId(ctx context.Context, roleId int32) ([]*models.Users2role, error) {

	var users2roles []*models.Users2role

	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleId).
		Preload("User").
		Preload("User.Profile").
		Find(&users2roles).Error
	if err != nil {
		return nil, err
	}

	return users2roles, nil

}

func (r *users2roleRepository) GetAllUsers2roles(ctx context.Context) ([]*models.Users2role, error) {

	var users2roles []*models.Users2role

	err := r.db.WithContext(ctx).Find(&users2roles).Error
	if err != nil {
		return nil, err
	}

	return users2roles, nil

}

func (r *users2roleRepository) PageUsers2role(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var users2roles []models.Users2role

	var err = r.db.WithContext(ctx).Scopes(pkg.Paginate(users2roles, &pagination, r.db)).Find(&users2roles).Error
	pagination.Rows = make([]any, len(users2roles))
	for i, users2role := range users2roles {
		pagination.Rows[i] = users2role
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
