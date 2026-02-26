package geofences

import (
	"context"
	"encoding/json"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type geofenceRepository struct{ db *gorm.DB }

func NewGeofenceRepository(db *gorm.DB) GeofenceRepository {
	return &geofenceRepository{db}
}

var _ GeofenceRepository = &geofenceRepository{}

func (r *geofenceRepository) CreateGeofence(ctx context.Context, input *models.CreateGeofenceInput, createdBy int) (*GeofenceDB, error) {
	// Konversi Any (interface{}) → GeoJSONCoords
	coords, err := toCoords(input.Coordinates)
	if err != nil {
		return nil, err
	}

	g := &GeofenceDB{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Color:       strOr(input.Color, "#3B82F6"),
		FillColor:   strOr(input.FillColor, "#3B82F640"),
		StrokeWidth: int32Or(input.StrokeWidth, 2),
		GeoType:     input.GeoType,
		Coordinates: coords,
		Radius:      input.Radius,
		IsActive:    boolOr(input.IsActive, true),
		CreatedBy:   createdBy,
	}
	err = r.db.WithContext(ctx).Create(g).Error
	return g, err
}

func (r *geofenceRepository) UpdateGeofence(ctx context.Context, id int32, input *models.UpdateGeofenceInput) (*GeofenceDB, error) {
	g := &GeofenceDB{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(g).Error; err != nil {
		return nil, err
	}
	applyUpdate(g, input)
	err := r.db.WithContext(ctx).Where("id = ?", id).Model(g).Updates(g).Error
	return g, err
}

func (r *geofenceRepository) UpdateGeofenceByUUID(ctx context.Context, uuid string, input *models.UpdateGeofenceInput) (*GeofenceDB, error) {
	g := &GeofenceDB{}
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(g).Error; err != nil {
		return nil, err
	}
	applyUpdate(g, input)
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(g).Updates(g).Error
	return g, err
}

func (r *geofenceRepository) DeleteGeofence(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&GeofenceDB{}).Error
}

func (r *geofenceRepository) DeleteGeofenceByUUID(ctx context.Context, uuid string) error {
	return r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&GeofenceDB{}).Error
}

func (r *geofenceRepository) GetGeofenceByID(ctx context.Context, id int32) (*GeofenceDB, error) {
	g := &GeofenceDB{}
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = 0", id).Take(g).Error
	return g, err
}

func (r *geofenceRepository) GetGeofenceByUUID(ctx context.Context, uuid string) (*GeofenceDB, error) {
	g := &GeofenceDB{}
	err := r.db.WithContext(ctx).Where("uuid = ? AND deleted_at = 0", uuid).Take(g).Error
	return g, err
}

func (r *geofenceRepository) GetAllGeofences(ctx context.Context) ([]*GeofenceDB, error) {
	var list []*GeofenceDB
	err := r.db.WithContext(ctx).Where("deleted_at = 0").Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *geofenceRepository) PageGeofence(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var list []GeofenceDB
	err := r.db.Scopes(pkg.Paginate(list, &pagination, r.db)).Find(&list).Error
	pagination.Rows = make([]any, len(list))
	for i, g := range list {
		pagination.Rows[i] = g
	}
	return pagination, err
}

// ─── helpers ───────────────────────────────────────────────────

func toCoords(v interface{}) (GeoJSONCoords, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var coords GeoJSONCoords
	return coords, json.Unmarshal(b, &coords)
}

func applyUpdate(g *GeofenceDB, input *models.UpdateGeofenceInput) {
	if input.Name != nil {
		g.Name = *input.Name
	}
	if input.Description != nil {
		g.Description = input.Description
	}
	if input.Type != nil {
		g.Type = *input.Type
	}
	if input.Color != nil {
		g.Color = *input.Color
	}
	if input.FillColor != nil {
		g.FillColor = *input.FillColor
	}
	if input.StrokeWidth != nil {
		g.StrokeWidth = int32(*input.StrokeWidth)
	}
	if input.GeoType != nil {
		g.GeoType = *input.GeoType
	}
	if input.Radius != nil {
		g.Radius = input.Radius
	}
	if input.IsActive != nil {
		g.IsActive = *input.IsActive
	}
	if input.Coordinates != nil {
		if coords, err := toCoords(input.Coordinates); err == nil {
			g.Coordinates = coords
		}
	}
}

func strOr(v *string, def string) string {
	if v != nil {
		return *v
	}
	return def
}

func int32Or(v *int, def int32) int32 {
	if v != nil {
		return int32(*v)
	}
	return def
}

func boolOr(v *bool, def bool) bool {
	if v != nil {
		return *v
	}
	return def
}
