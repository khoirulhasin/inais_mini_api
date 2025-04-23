package interfaces

import (
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/answers"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/question_options"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/questions"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AnsService            answers.AnsService
	QuestionService       questions.QuesService
	QuestionOptionService question_options.OptService
}
