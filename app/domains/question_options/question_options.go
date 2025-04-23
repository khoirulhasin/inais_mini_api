package question_options

import "github.com/khoirulhasin/globe_tracker_api/app/models"

type OptRepository interface {
	CreateQuestionOption(question *models.QuestionOption) (*models.QuestionOption, error)
	UpdateQuestionOption(question *models.QuestionOption) (*models.QuestionOption, error)
	DeleteQuestionOption(id string) error
	DeleteQuestionOptionByQuestionID(questionId string) error
	GetQuestionOptionByID(id string) (*models.QuestionOption, error)
	GetQuestionOptionByQuestionID(questionId string) ([]*models.QuestionOption, error)
}
