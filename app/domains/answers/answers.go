package answers

import (
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

type AnsRepository interface {
	CreateAnswer(answer *models.Answer) (*models.Answer, error)
	UpdateAnswer(answer *models.Answer) (*models.Answer, error)
	DeleteAnswer(id string) error
	GetAnswerByID(id string) (*models.Answer, error)
	GetAllQuestionAnswers(questionId string) ([]*models.Answer, error)
	GetQuestionOptionByID(id string) (*models.QuestionOption, error)
}
