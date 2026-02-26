package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/untirta_api/app/dependencies"
)

func SetupAllRoutes(r *gin.Engine, handlers *dependencies.Handlers) {
	api := r.Group("/api/v1")
	{
		// Setup marker routes
		SetupMarkerRoutes(api, handlers.MarkerHandler)

	}
}
