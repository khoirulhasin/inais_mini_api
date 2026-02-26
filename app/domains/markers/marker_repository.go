package markers

import (
	"context"
	"fmt"
	"time"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type markerRepository struct {
	db *gorm.DB
}

func NewMarkerRepository(db *gorm.DB) *markerRepository {
	return &markerRepository{
		db,
	}
}

var _ MarkerRepository = &markerRepository{}

func (r *markerRepository) CreateMarker(ctx context.Context, marker *models.Marker) (*models.Marker, error) {

	err := r.db.WithContext(ctx).Create(&marker).Error
	if err != nil {
		return nil, err
	}

	return marker, nil
}

func (r *markerRepository) UpdateMarker(ctx context.Context, id int32, marker *models.Marker) (*models.Marker, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&marker).Updates(marker).Error
	if err != nil {
		return nil, err
	}

	return marker, nil
}

func (r *markerRepository) UpdateMarkerByUUID(ctx context.Context, uuid string, marker *models.Marker) (*models.Marker, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&marker).Updates(marker).Error
	if err != nil {
		return nil, err
	}

	return marker, nil
}

func (s *markerRepository) DeleteMarker(ctx context.Context, id int32) error {

	marker := &models.Marker{}

	err := s.db.WithContext(ctx).Where("id = ?", id).Delete(marker).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *markerRepository) DeleteMarkerByUUID(ctx context.Context, uuid string) error {

	marker := &models.Marker{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(marker).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *markerRepository) GetMarkerByID(ctx context.Context, id int32) (*models.Marker, error) {

	var marker = &models.Marker{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&marker).Error
	if err != nil {
		return nil, err
	}

	return marker, nil
}

func (r *markerRepository) GetMarkerByUUID(ctx context.Context, uuid string) (*models.Marker, error) {

	var marker = &models.Marker{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&marker).Error
	if err != nil {
		return nil, err
	}

	return marker, nil
}

func (r *markerRepository) GetAllMarkers(ctx context.Context) ([]*models.Marker, error) {

	var markers []*models.Marker

	err := r.db.WithContext(ctx).Preload("MarkerType").Find(&markers).Error
	if err != nil {
		return nil, err
	}

	return markers, nil

}

func (r *markerRepository) PageMarker(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var markers []models.Marker

	var err = r.db.Scopes(pkg.Paginate(markers, &pagination, r.db)).Find(&markers).Error
	pagination.Rows = make([]any, len(markers))
	for i, marker := range markers {
		pagination.Rows[i] = marker
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}

func (r *markerRepository) GetNearestMarkers(ctx context.Context, imei string, lat, lng float64, limit int) ([]*NearestMarkerResponse, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}

	currentTime := time.Now().Unix()

	distanceFormula := fmt.Sprintf(
		"(6371 * acos(cos(radians(%f)) * cos(radians(lat)) * cos(radians(lng) - radians(%f)) + sin(radians(%f)) * sin(radians(lat))))",
		lat, lng, lat,
	)

	query := fmt.Sprintf(`
		SELECT 
			m.lat, 
			m.lng, 
			%s as distance,
			mt.id as marker_type_id,
			mt.name as marker_type_name,
			mt.icon as marker_type_icon
		FROM markers m
		LEFT JOIN marker_types mt ON m.marker_type_id = mt.id
		WHERE (m.start IS NULL OR m.start <= ?) 
		AND (m."end" IS NULL OR m."end" >= ?) 
		AND m.deleted_at = 0
		ORDER BY distance ASC 
		LIMIT ?`, distanceFormula)

	var results []struct {
		Lat            float64 `gorm:"column:lat"`
		Lng            float64 `gorm:"column:lng"`
		Distance       float64 `gorm:"column:distance"`
		MarkerTypeID   int     `gorm:"column:marker_type_id"`
		MarkerTypeName string  `gorm:"column:marker_type_name"`
		MarkerTypeIcon string  `gorm:"column:marker_type_icon"`
	}

	err := r.db.WithContext(ctx).Raw(query, currentTime, currentTime, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var response []*NearestMarkerResponse
	for _, result := range results {
		response = append(response, &NearestMarkerResponse{
			Lat:      result.Lat,
			Lng:      result.Lng,
			Distance: result.Distance,
			MarkerType: MarkerTypeSimple{
				ID:   result.MarkerTypeID,
				Name: result.MarkerTypeName,
			},
		})
	}

	return response, nil
}
