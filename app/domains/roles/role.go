package roles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	UpdateRole(ctx context.Context, id int32, role *models.Role) (*models.Role, error)
	UpdateRoleByUUID(ctx context.Context, uuid string, role *models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, id int32) error
	DeleteRoleByUUID(ctx context.Context, uuid string) error
	GetRoleByID(ctx context.Context, id int32) (*models.Role, error)
	GetRoleByUUID(ctx context.Context, uuid string) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]*models.Role, error)
	PageRole(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
