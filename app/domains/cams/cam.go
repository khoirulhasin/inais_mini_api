package cams

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type CamRepository interface {
	CreateCam(ctx context.Context, cam *models.Cam) (*models.Cam, error)
	UpdateCam(ctx context.Context, id int32, cam *models.Cam) (*models.Cam, error)
	UpdateCamByUUID(ctx context.Context, uuid string, cam *models.Cam) (*models.Cam, error)
	DeleteCam(ctx context.Context, id int32) error
	DeleteCamByUUID(ctx context.Context, uuid string) error
	GetCamByID(ctx context.Context, id int32) (*models.Cam, error)
	GetCamByUUID(ctx context.Context, uuid string) (*models.Cam, error)
	GetCamByStateID(ctx context.Context, stateId int32) ([]*models.Cam, error)
	GetAllCams(ctx context.Context) ([]*models.Cam, error)
	PageCam(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
