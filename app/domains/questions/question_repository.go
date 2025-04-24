package questions

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

type quesRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *quesRepository {
	return &quesRepository{
		db,
	}
}

// We implement the interface defined in the domain
var _ QuesRepository = &quesRepository{}

func (s *quesRepository) CreateQuestion(question *models.Question) (*models.Question, error) {

	err := s.db.Create(&question).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("question title already taken")
		}
		return nil, err
	}

	return question, nil

}

func (s *quesRepository) UpdateQuestion(question *models.Question) (*models.Question, error) {

	err := s.db.Save(&question).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("question title already taken")
		}
		return nil, err
	}

	return question, nil

}

func (s *quesRepository) DeleteQuestion(id string) error {

	ques := &models.Question{}

	err := s.db.Where("id = ?", id).Delete(&ques).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *quesRepository) GetQuestionByID(id string) (*models.Question, error) {

	ques := &models.Question{}

	err := s.db.Where("id = ?", id).Preload("QuestionOption").Take(&ques).Error
	if err != nil {
		return nil, err
	}

	return ques, nil
}

func (s *quesRepository) GetAllQuestions() ([]*models.Question, error) {

	var questions []*models.Question

	err := s.db.Preload("QuestionOption").Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}
