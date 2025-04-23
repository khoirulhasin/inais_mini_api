package questions

import "github.com/khoirulhasin/globe_tracker_api/app/models"

type QuesService interface {
	CreateQuestion(question *models.Question) (*models.Question, error)
	UpdateQuestion(question *models.Question) (*models.Question, error)
	DeleteQuestion(id string) error
	GetQuestionByID(id string) (*models.Question, error)
	GetAllQuestions() ([]*models.Question, error)
}
