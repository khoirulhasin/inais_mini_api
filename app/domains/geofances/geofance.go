package geofences

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/plugin/soft_delete"
)

// GeofenceDB adalah struct GORM — ditulis manual karena perlu
// custom type GeoJSONCoords untuk JSONB
type GeofenceDB struct {
	ID          int32                 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UUID        string                `json:"uuid" gorm:"column:uuid;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	Name        string                `json:"name" gorm:"column:name;not null"`
	Description *string               `json:"description" gorm:"column:description"`
	Type        string                `json:"type" gorm:"column:type;not null"`
	Color       string                `json:"color" gorm:"column:color;default:'#3B82F6'"`
	FillColor   string                `json:"fillColor" gorm:"column:fill_color;default:'#3B82F640'"`
	StrokeWidth int32                 `json:"strokeWidth" gorm:"column:stroke_width;default:2"`
	GeoType     string                `json:"geoType" gorm:"column:geo_type;not null"`
	Coordinates GeoJSONCoords         `json:"coordinates" gorm:"column:coordinates;type:jsonb;not null"`
	Radius      *float64              `json:"radius" gorm:"column:radius"`
	IsActive    bool                  `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedBy   int                   `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy   *int                  `json:"updatedBy" gorm:"column:updated_by"`
	CreatedAt   int64                 `json:"createdAt" gorm:"column:created_at;type:bigint;autoCreateTime:milli"`
	UpdatedAt   int64                 `json:"updatedAt" gorm:"column:updated_at;type:bigint;autoUpdateTime:milli"`
	DeletedAt   soft_delete.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;type:bigint;softDelete:milli;default:0"`
}

func (GeofenceDB) TableName() string { return "geofences" }

type GeofenceRepository interface {
	CreateGeofence(ctx context.Context, input *models.CreateGeofenceInput, createdBy int) (*GeofenceDB, error)
	UpdateGeofence(ctx context.Context, id int32, input *models.UpdateGeofenceInput) (*GeofenceDB, error)
	UpdateGeofenceByUUID(ctx context.Context, uuid string, input *models.UpdateGeofenceInput) (*GeofenceDB, error)
	DeleteGeofence(ctx context.Context, id int32) error
	DeleteGeofenceByUUID(ctx context.Context, uuid string) error
	GetGeofenceByID(ctx context.Context, id int32) (*GeofenceDB, error)
	GetGeofenceByUUID(ctx context.Context, uuid string) (*GeofenceDB, error)
	GetAllGeofences(ctx context.Context) ([]*GeofenceDB, error)
	PageGeofence(ctx context.Context, pagination models.Pagination) (models.Pagination, error)
}
