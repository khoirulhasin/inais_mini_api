package users

import (
	"context"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/configs"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/helpers"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db,
	}
}

var _ UserRepository = &userRepository{}

func (r *userRepository) Login(ctx context.Context, account, password string) (*AuthenticationResponse, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", account, account).Take(&user).Error
	if err != nil {
		return nil, err
	}

	err = helpers.VerifyPassword(password, user.Password)

	if err != nil {
		return nil, err
	}
	users2roles := []*models.Users2role{}

	err = r.db.WithContext(ctx).Where("user_id = ?", user.ID).Preload("User").Preload("Role").Find(&users2roles).Error

	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(users2roles))
	for _, users2role := range users2roles {
		if users2role.Role != nil {
			roles = append(roles, users2role.Role.Code)
		}
	}

	claims := CustomClaims{
		ID:               user.ID,
		UUID:             user.UUID.String(),
		Username:         user.Username,
		Email:            user.Email,
		Roles:            roles,
		RegisteredClaims: configs.GetClaim(user.ID),
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	return &AuthenticationResponse{
		Token:    tokenString,
		ID:       user.ID,
		UUID:     user.UUID.String(),
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int32, user *models.User) (*models.User, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&user).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserByUUID(ctx context.Context, uuid string, user *models.User) (*models.User, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&user).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int32, password string) (*models.User, error) {

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Model(&models.User{}).Updates(models.User{Password: password}).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *userRepository) UpdatePasswordByUUID(ctx context.Context, uuid string, password string) (*models.User, error) {

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Model(&models.User{}).Updates(models.User{Password: password}).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int32) error {

	var user models.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}

	return nil

}

func (r *userRepository) DeleteUserByUUID(ctx context.Context, uuid string) error {

	user := &models.User{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(user).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (*models.User, error) {

	var user = &models.User{}

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Role").
		Preload("Profile").
		Take(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {

	var user = &models.User{}

	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Profile").
		Where("users.uuid = ?", uuid).
		Take(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {

	var users []*models.User

	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (r *userRepository) PageUser(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var users []models.User

	var err = r.db.WithContext(ctx).Preload("Role").Scopes(pkg.Paginate(models.User{}, &pagination, r.db)).Find(&users).Error
	pagination.Rows = make([]any, len(users))
	for i, user := range users {
		pagination.Rows[i] = user
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}

func (r *userRepository) PageUserByRoleIDs(ctx context.Context, roleIds []int32, pagination models.Pagination) (models.Pagination, error) {
	var users []models.User

	var err = r.db.WithContext(ctx).
		Preload("Role").
		Preload("Profile").
		Joins("LEFT JOIN users2roles ON users.id = users2roles.user_id").
		Joins("LEFT JOIN profiles ON users.id = profiles.user_id").
		Where("users2roles.role_id IN (?)", roleIds).
		Scopes(pkg.Paginate(models.User{}, &pagination, r.db)).
		Find(&users).Error
	pagination.Rows = make([]any, len(users))
	for i, user := range users {
		pagination.Rows[i] = user
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
