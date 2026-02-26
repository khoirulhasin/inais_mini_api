package marker_types

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type MarkerTypeRepository interface {
	CreateMarkerType(ctx context.Context, markerType *models.MarkerType) (*models.MarkerType, error)
	UpdateMarkerType(ctx context.Context, id int32, markerType *models.MarkerType) (*models.MarkerType, error)
	UpdateMarkerTypeByUUID(ctx context.Context, uuid string, markerType *models.MarkerType) (*models.MarkerType, error)
	DeleteMarkerType(ctx context.Context, id int32) error
	DeleteMarkerTypeByUUID(ctx context.Context, uuid string) error
	GetMarkerTypeByID(ctx context.Context, id int32) (*models.MarkerType, error)
	GetMarkerTypeByUUID(ctx context.Context, uuid string) (*models.MarkerType, error)
	GetAllMarkerTypes(ctx context.Context) ([]*models.MarkerType, error)
	PageMarkerType(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
