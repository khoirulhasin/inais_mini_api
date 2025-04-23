package dependencies

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/answers"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/question_options"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/questions"
	"github.com/khoirulhasin/globe_tracker_api/app/generated"
	"github.com/khoirulhasin/globe_tracker_api/app/infrastructures/dbs/postgres"
	"github.com/khoirulhasin/globe_tracker_api/app/interfaces"
	"github.com/vektah/gqlparser/v2/ast"
)

func Init() http.Handler {

	var connPostgres = postgres.Connect()

	var ansRepository answers.AnsRepository
	var questionRepository questions.QuesRepository
	var questionOptRepository question_options.OptRepository

	ansRepository = answers.NewAnswerRepository(connPostgres)
	questionRepository = questions.NewQuestionRepository(connPostgres)
	questionOptRepository = question_options.NewQuestionOptionRepository(connPostgres)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &interfaces.Resolver{
		AnsRepository:            ansRepository,
		QuestionRepository:       questionRepository,
		QuestionOptionRepository: questionOptRepository,
	}}))

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return h
}
