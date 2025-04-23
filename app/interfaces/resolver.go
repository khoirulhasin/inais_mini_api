package interfaces

import (
	"github.com/khoirulhasin/globe_tracker_api/app/domains/answers"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/question_options"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/questions"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AnsRepository            answers.AnsRepository
	QuestionRepository       questions.QuesRepository
	QuestionOptionRepository question_options.OptRepository
}
