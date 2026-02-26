package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/untirta_api/app/api/handlers"
)

func SetupMarkerRoutes(api *gin.RouterGroup, markerHandler *handlers.MarkerHandler) {
	markers := api.Group("/markers")
	{
		markers.GET("", markerHandler.GetMarkersWithPagination)
		markers.GET("/all", markerHandler.GetAllMarkers)
		markers.GET("/nearest", markerHandler.GetNearestMarkers)
		markers.GET("/:id", markerHandler.GetMarkerByID)
		markers.GET("/uuid/:uuid", markerHandler.GetMarkerByUUID)
		markers.POST("", markerHandler.CreateMarker)
		markers.PUT("/:id", markerHandler.UpdateMarker)
		markers.PUT("/uuid/:uuid", markerHandler.UpdateMarkerByUUID)
		markers.DELETE("/:id", markerHandler.DeleteMarker)
		markers.DELETE("/uuid/:uuid", markerHandler.DeleteMarkerByUUID)
	}
}
