package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khoirulhasin/globe_tracker_api/app/dependencies"
	"github.com/khoirulhasin/globe_tracker_api/app/infrastructures/middlewares"
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
	r.Use(middlewares.GinContextToContextMiddleware())

	r.POST("/query", dependencies.GraphqlHandler())
	r.GET("/", dependencies.PlaygroundHandler())
	r.Run()
}
