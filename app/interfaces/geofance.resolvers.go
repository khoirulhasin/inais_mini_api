package interfaces

// File ini di-generate stub-nya oleh gqlgen, isi implementasinya manual.

import (
	"context"

	"github.com/google/uuid"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/helpers"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
)

func (r *queryResolver) GetAllGeofences(ctx context.Context) (interface{}, error) {
	return r.GeofenceRepository.GetAllGeofences(ctx)
}

func (r *queryResolver) GetOneGeofence(ctx context.Context, id int) (interface{}, error) {
	return r.GeofenceRepository.GetGeofenceByID(ctx, int32(id))
}

func (r *queryResolver) GetOneGeofenceByUUID(ctx context.Context, uuid uuid.UUID) (interface{}, error) {
	return r.GeofenceRepository.GetGeofenceByUUID(ctx, uuid.String())
}

func (r *queryResolver) PageGeofence(ctx context.Context, pageInput *models.PageInput) (interface{}, error) {
	limit, offset, sortField, sortOrder, search, filters := pkg.PageInputIsNil(pageInput)

	// pagination := models.Pagination{
	// 	Limit: limit, Offset: offset,
	// 	SortField: sortField, SortOrder: sortOrder,
	// 	Search: search, Filters: filters,
	// }

	var mappedFilters []*models.Filter
	for _, f := range filters {
		mappedFilters = append(mappedFilters, &models.Filter{
			Key:      f.Key,
			Operator: f.Operator,
			Value:    f.Value,
		})
	}

	pagination := models.Pagination{
		Limit:     &limit,
		Offset:    &offset,
		SortField: &sortField,
		SortOrder: &sortOrder,
		Search:    &search,
		Filters:   mappedFilters,
	}

	return r.GeofenceRepository.PageGeofence(ctx, pagination)
}

func (r *mutationResolver) CreateGeofence(ctx context.Context, input models.CreateGeofenceInput) (interface{}, error) {
	token, err := helpers.GetToken(ctx)

	if err != nil {
		return nil, err
	}

	userID, err := helpers.GetUserID(token.(string))

	if err != nil {
		return nil, err
	}
	return r.GeofenceRepository.CreateGeofence(ctx, &input, userID)
}

func (r *mutationResolver) UpdateGeofence(ctx context.Context, id int, input models.UpdateGeofenceInput) (interface{}, error) {
	return r.GeofenceRepository.UpdateGeofence(ctx, int32(id), &input)
}

func (r *mutationResolver) UpdateGeofenceByUUID(ctx context.Context, uuid uuid.UUID, input models.UpdateGeofenceInput) (interface{}, error) {
	return r.GeofenceRepository.UpdateGeofenceByUUID(ctx, uuid.String(), &input)
}

func (r *mutationResolver) DeleteGeofence(ctx context.Context, id int) (interface{}, error) {
	err := r.GeofenceRepository.DeleteGeofence(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// DeleteGeofenceByUUID is the resolver for the DeleteGeofenceByUuid field.
func (r *mutationResolver) DeleteGeofenceByUUID(ctx context.Context, uuid uuid.UUID) (interface{}, error) {
	err := r.GeofenceRepository.DeleteGeofenceByUUID(ctx, uuid.String())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
