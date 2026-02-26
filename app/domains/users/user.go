package users

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khoirulhasin/untirta_api/app/models"
)

type UserRepository interface {
	Login(ctx context.Context, account, password string) (*AuthenticationResponse, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, id int32, user *models.User) (*models.User, error)
	UpdateUserByUUID(ctx context.Context, uuid string, user *models.User) (*models.User, error)
	UpdatePassword(ctx context.Context, id int32, password string) (*models.User, error)
	UpdatePasswordByUUID(ctx context.Context, uuid string, password string) (*models.User, error)
	DeleteUser(ctx context.Context, id int32) error
	DeleteUserByUUID(ctx context.Context, uuid string) error
	GetUserByID(ctx context.Context, id int32) (*models.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	PageUser(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
	PageUserByRoleIDs(ctx context.Context, roleIds []int32, pagination models.Pagination) (models.Pagination, error)
}

type CustomClaims struct {
	ID       int
	UUID     string
	Username string
	Email    string
	Roles    []string
	jwt.RegisteredClaims
}

type AuthenticationResponse struct {
	Token    string   `json:"token"`
	ID       int      `json:"id"`
	UUID     string   `json:"uuid"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}
