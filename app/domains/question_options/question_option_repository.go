package question_options

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

type optRepository struct {
	db *gorm.DB
}

func NewQuestionOptionRepository(db *gorm.DB) *optRepository {
	return &optRepository{
		db,
	}
}

// We implement the interface defined in the domain
var _ OptRepository = &optRepository{}

func (s *optRepository) CreateQuestionOption(questOpt *models.QuestionOption) (*models.QuestionOption, error) {

	//check if this question option title or the position or the correctness already exist for the question
	oldOpts, _ := s.GetQuestionOptionByQuestionID(questOpt.QuestionID)
	if len(oldOpts) > 0 {
		for _, v := range oldOpts {
			if v.Title == questOpt.Title || v.Position == questOpt.Position || (v.IsCorrect == true && questOpt.IsCorrect == true) {
				return nil, errors.New("two question options can't have the same title, position and/or the same correct answer")
			}
		}
	}

	err := s.db.Create(&questOpt).Error
	if err != nil {
		return nil, err
	}

	return questOpt, nil
}

func (s *optRepository) UpdateQuestionOption(questOpt *models.QuestionOption) (*models.QuestionOption, error) {

	err := s.db.Save(&questOpt).Error
	if err != nil {
		return nil, err
	}

	return questOpt, nil

}

func (s *optRepository) DeleteQuestionOption(id string) error {

	err := s.db.Delete(id).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *optRepository) DeleteQuestionOptionByQuestionID(questId string) error {

	err := s.db.Delete(questId).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *optRepository) GetQuestionOptionByID(id string) (*models.QuestionOption, error) {

	quesOpt := &models.QuestionOption{}

	err := s.db.Where("id = ?", id).Take(&quesOpt).Error
	if err != nil {
		return nil, err
	}

	return quesOpt, nil

}

func (s *optRepository) GetQuestionOptionByQuestionID(id string) ([]*models.QuestionOption, error) {

	var quesOpts []*models.QuestionOption

	err := s.db.Where("question_id = ?", id).Find(&quesOpts).Error
	if err != nil {
		return nil, err
	}

	return quesOpts, nil
}
