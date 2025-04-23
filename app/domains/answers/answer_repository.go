package answers

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

type ansRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *ansRepository {
	return &ansRepository{
		db,
	}
}

// We implement the interface defined in the domain
var _ AnsRepository = &ansRepository{}

func (s *ansRepository) CreateAnswer(answer *models.Answer) (*models.Answer, error) {

	//first we need to check if the ans have been entered for this question:
	oldAns, _ := s.GetAllQuestionAnswers(answer.QuestionID)
	if len(oldAns) > 0 {
		for _, v := range oldAns {
			//We cannot have two correct answers for this type of quiz
			if v.IsCorrect == true && answer.IsCorrect {
				return nil, errors.New("cannot have two correct answers for the same question")
			}
		}
	}

	err := s.db.Create(&answer).Error
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (s *ansRepository) UpdateAnswer(answer *models.Answer) (*models.Answer, error) {

	err := s.db.Save(&answer).Error
	if err != nil {
		return nil, err
	}

	return answer, nil

}

func (s *ansRepository) DeleteAnswer(id string) error {

	ans := &models.Answer{}

	err := s.db.Where("id = ?", id).Delete(ans).Error
	if err != nil {
		return err
	}

	return nil

}

func (s *ansRepository) GetAnswerByID(id string) (*models.Answer, error) {

	var ans = &models.Answer{}

	err := s.db.Where("id = ?", id).Take(&ans).Error
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (s *ansRepository) GetAllQuestionAnswers(questionId string) ([]*models.Answer, error) {

	var answers []*models.Answer

	err := s.db.Where("question_id = ?", questionId).Find(&answers).Error
	if err != nil {
		return nil, err
	}

	return answers, nil

}

func (s *ansRepository) GetQuestionOptionByID(id string) (*models.QuestionOption, error) {

	quesOpt := &models.QuestionOption{}

	err := s.db.Where("id = ?", id).Take(&quesOpt).Error
	if err != nil {
		return nil, err
	}

	return quesOpt, nil

}
