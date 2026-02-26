package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khoirulhasin/untirta_api/app/api/routes"
	"github.com/khoirulhasin/untirta_api/app/dependencies"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/middlewares"
	_ "github.com/lib/pq"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	r := gin.Default()

	// Initialize GraphQL first untuk trigger dependency injection
	// _ = dependencies.GraphqlHandler(r)

	// GraphQL endpoints
	r.POST("/query", middlewares.HeaderToContextMiddleware(), dependencies.GraphqlHandler(r))
	r.GET("/", dependencies.PlaygroundHandler())

	// Setup REST API routes
	handlers := dependencies.GetHandlers()
	if handlers != nil {
		routes.SetupAllRoutes(r, handlers)
	}

	r.Run()
}
