package profiles

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, country *models.Profile) (*models.Profile, error)
	UpdateProfile(ctx context.Context, id int32, country *models.Profile) (*models.Profile, error)
	UpdateProfileByUUID(ctx context.Context, uuid string, country *models.Profile) (*models.Profile, error)
	UpdateProfileByUserID(ctx context.Context, userId int32, profile *models.Profile) (*models.Profile, error)
	DeleteProfile(ctx context.Context, id int32) error
	DeleteProfileByUUID(ctx context.Context, uuid string) error
	GetProfileByID(ctx context.Context, id int32) (*models.Profile, error)
	GetProfileByUUID(ctx context.Context, uuid string) (*models.Profile, error)
	GetProfileByUserID(ctx context.Context, userId int32) (*models.Profile, error)
	GetAllProfiles(ctx context.Context) ([]*models.Profile, error)
	PageProfile(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
