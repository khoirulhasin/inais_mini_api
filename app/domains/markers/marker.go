package markers

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
)

type MarkerRepository interface {
	CreateMarker(ctx context.Context, marker *models.Marker) (*models.Marker, error)
	UpdateMarker(ctx context.Context, id int32, marker *models.Marker) (*models.Marker, error)
	UpdateMarkerByUUID(ctx context.Context, uuid string, marker *models.Marker) (*models.Marker, error)
	DeleteMarker(ctx context.Context, id int32) error
	DeleteMarkerByUUID(ctx context.Context, uuid string) error
	GetMarkerByID(ctx context.Context, id int32) (*models.Marker, error)
	GetMarkerByUUID(ctx context.Context, uuid string) (*models.Marker, error)
	GetAllMarkers(ctx context.Context) ([]*models.Marker, error)
	PageMarker(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
	GetNearestMarkers(ctx context.Context, imei string, lat, lng float64, limit int) ([]*NearestMarkerResponse, error)
}

type NearestMarkerResponse struct {
	Lat        float64          `json:"lat"`
	Lng        float64          `json:"lng"`
	Distance   float64          `json:"distance"`
	MarkerType MarkerTypeSimple `json:"markerType"`
}

type MarkerTypeSimple struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
