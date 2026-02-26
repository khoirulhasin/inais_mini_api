package geofences

import (
	"context"

	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/plugin/soft_delete"
)

// GeofenceDB adalah struct GORM — ditulis manual karena perlu
// custom type GeoJSONCoords untuk JSONB
type GeofenceDB struct {
	ID          int32                 `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        string                `gorm:"column:uuid;uniqueIndex;type:uuid;default:uuid_generate_v4()"`
	Name        string                `gorm:"column:name;not null"`
	Description *string               `gorm:"column:description"`
	Type        string                `gorm:"column:type;not null"`
	Color       string                `gorm:"column:color;default:'#3B82F6'"`
	FillColor   string                `gorm:"column:fill_color;default:'#3B82F640'"`
	StrokeWidth int32                 `gorm:"column:stroke_width;default:2"`
	GeoType     string                `gorm:"column:geo_type;not null"`
	Coordinates GeoJSONCoords         `gorm:"column:coordinates;type:jsonb;not null"`
	Radius      *float64              `gorm:"column:radius"`
	IsActive    bool                  `gorm:"column:is_active;default:true"`
	CreatedBy   int                   `gorm:"column:created_by"`
	UpdatedBy   *int                  `gorm:"column:updated_by"`
	CreatedAt   int64                 `gorm:"column:created_at;type:bigint;autoCreateTime:milli"`
	UpdatedAt   int64                 `gorm:"column:updated_at;type:bigint;autoUpdateTime:milli"`
	DeletedAt   soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint;softDelete:milli;default:0"`
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
