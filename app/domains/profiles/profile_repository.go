package profiles

import (
	"context"
	"errors"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepository {
	return &profileRepository{
		db,
	}
}

var _ ProfileRepository = &profileRepository{}

func (r *profileRepository) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {

	err := r.db.WithContext(ctx).Create(&profile).Error
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, id int32, profile *models.Profile) (*models.Profile, error) {

	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&profile).Updates(profile).Error
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) UpdateProfileByUUID(ctx context.Context, uuid string, profile *models.Profile) (*models.Profile, error) {

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Model(&profile).Updates(profile).Error
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) UpdateProfileByUserID(ctx context.Context, userId int32, profile *models.Profile) (*models.Profile, error) {

	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Model(&profile).Updates(profile).Error
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) DeleteProfile(ctx context.Context, id int32) error {

	profile := &models.Profile{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(profile).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *profileRepository) DeleteProfileByUUID(ctx context.Context, uuid string) error {

	profile := &models.Profile{}

	err := s.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(profile).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *profileRepository) GetProfileByID(ctx context.Context, id int32) (*models.Profile, error) {

	var profile = &models.Profile{}

	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found, return empty slice instead of error
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) GetProfileByUUID(ctx context.Context, uuid string) (*models.Profile, error) {

	var profile = &models.Profile{}

	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).Take(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found, return empty slice instead of error
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) GetProfileByUserID(ctx context.Context, userId int32) (*models.Profile, error) {

	var profile = &models.Profile{}

	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Take(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found, return empty slice instead of error
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}

func (r *profileRepository) GetAllProfiles(ctx context.Context) ([]*models.Profile, error) {

	var profiles []*models.Profile

	err := r.db.WithContext(ctx).Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	return profiles, nil

}

func (r *profileRepository) PageProfile(ctx context.Context, pagination models.Pagination) (models.Pagination, error) {
	var profiles []models.Profile

	var err = r.db.WithContext(ctx).Scopes(pkg.Paginate(profiles, &pagination, r.db)).Find(&profiles).Error
	pagination.Rows = make([]any, len(profiles))
	for i, profile := range profiles {
		pagination.Rows[i] = profile
	}

	if err != nil {
		return pagination, err
	}

	return pagination, nil

}
