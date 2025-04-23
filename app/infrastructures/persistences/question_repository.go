package persistences

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/questions"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

type quesService struct {
	db *gorm.DB
}

func NewQuestion(db *gorm.DB) *quesService {
	return &quesService{
		db,
	}
}

// We implement the interface defined in the domain
var _ questions.QuesService = &quesService{}

func (s *quesService) CreateQuestion(question *models.Question) (*models.Question, error) {

	err := s.db.Create(&question).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("question title already taken")
		}
		return nil, err
	}

	return question, nil

}

func (s *quesService) UpdateQuestion(question *models.Question) (*models.Question, error) {

	err := s.db.Save(&question).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("question title already taken")
		}
		return nil, err
	}

	return question, nil

}

func (s *quesService) DeleteQuestion(id string) error {

	ques := &models.Question{}

	err := s.db.Where("id = ?", id).Delete(&ques).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *quesService) GetQuestionByID(id string) (*models.Question, error) {

	ques := &models.Question{}

	err := s.db.Where("id = ?", id).Preload("QuestionOption").Take(&ques).Error
	if err != nil {
		return nil, err
	}

	return ques, nil
}

func (s *quesService) GetAllQuestions() ([]*models.Question, error) {

	var questions []*models.Question

	err := s.db.Preload("QuestionOption").Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}
