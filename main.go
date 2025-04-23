package main

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/answers"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/question_options"
	"github.com/khoirulhasin/globe_tracker_api/app/domains/repositories/questions"
	"github.com/khoirulhasin/globe_tracker_api/app/generated"
	"github.com/khoirulhasin/globe_tracker_api/app/infrastructures/persistences"
	"github.com/khoirulhasin/globe_tracker_api/app/infrastructures/persistences/dbs"
	"github.com/khoirulhasin/globe_tracker_api/app/interfaces"
	_ "github.com/lib/pq"

	"net/http"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	var (
		defaultPort      = "8080"
		databaseUser     = os.Getenv("DATABASE_USER")
		databaseName     = os.Getenv("DATABASE_NAME")
		databaseHost     = os.Getenv("DATABASE_HOST")
		databasePort     = os.Getenv("DATABASE_PORT")
		databasePassword = os.Getenv("DATABASE_PASSWORD")
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	connPostgres := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	conn := dbs.OpenPostgres(connPostgres)

	var ansService answers.AnsService
	var questionService questions.QuesService
	var questionOptService question_options.OptService

	ansService = persistences.NewAnswer(conn)
	questionService = persistences.NewQuestion(conn)
	questionOptService = persistences.NewQuestionOption(conn)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &interfaces.Resolver{
		AnsService:            ansService,
		QuestionService:       questionService,
		QuestionOptionService: questionOptService,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
