package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/untirta_api/app/domains/markers"
	"github.com/khoirulhasin/untirta_api/app/models"
)

type MarkerHandler struct {
	markerRepo markers.MarkerRepository
}

func NewMarkerHandler(markerRepo markers.MarkerRepository) *MarkerHandler {
	return &MarkerHandler{
		markerRepo: markerRepo,
	}
}

// GetAllMarkers godoc
// @Summary Get all markers
// @Description Get all markers with optional preloading
// @Tags markers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers [get]
func (h *MarkerHandler) GetAllMarkers(c *gin.Context) {
	ctx := c.Request.Context()

	markers, err := h.markerRepo.GetAllMarkers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get markers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   markers,
		"count":  len(markers),
	})
}

// GetMarkerByID godoc
// @Summary Get marker by ID
// @Description Get a single marker by ID
// @Tags markers
// @Accept json
// @Produce json
// @Param id path int true "Marker ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/markers/{id} [get]
func (h *MarkerHandler) GetMarkerByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid marker ID",
			"error":   err.Error(),
		})
		return
	}

	marker, err := h.markerRepo.GetMarkerByID(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Marker not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   marker,
	})
}

// GetMarkerByUUID godoc
// @Summary Get marker by UUID
// @Description Get a single marker by UUID
// @Tags markers
// @Accept json
// @Produce json
// @Param uuid path string true "Marker UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/markers/uuid/{uuid} [get]
func (h *MarkerHandler) GetMarkerByUUID(c *gin.Context) {
	ctx := c.Request.Context()

	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "UUID is required",
		})
		return
	}

	marker, err := h.markerRepo.GetMarkerByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Marker not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   marker,
	})
}

// CreateMarker godoc
// @Summary Create a new marker
// @Description Create a new marker
// @Tags markers
// @Accept json
// @Produce json
// @Param marker body models.Marker true "Marker data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers [post]
func (h *MarkerHandler) CreateMarker(c *gin.Context) {
	ctx := c.Request.Context()

	var marker models.Marker
	if err := c.ShouldBindJSON(&marker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	createdMarker, err := h.markerRepo.CreateMarker(ctx, &marker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create marker",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Marker created successfully",
		"data":    createdMarker,
	})
}

// UpdateMarker godoc
// @Summary Update marker by ID
// @Description Update an existing marker by ID
// @Tags markers
// @Accept json
// @Produce json
// @Param id path int true "Marker ID"
// @Param marker body models.Marker true "Marker data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers/{id} [put]
func (h *MarkerHandler) UpdateMarker(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid marker ID",
			"error":   err.Error(),
		})
		return
	}

	var marker models.Marker
	if err := c.ShouldBindJSON(&marker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	updatedMarker, err := h.markerRepo.UpdateMarker(ctx, int32(id), &marker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update marker",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Marker updated successfully",
		"data":    updatedMarker,
	})
}

// UpdateMarkerByUUID godoc
// @Summary Update marker by UUID
// @Description Update an existing marker by UUID
// @Tags markers
// @Accept json
// @Produce json
// @Param uuid path string true "Marker UUID"
// @Param marker body models.Marker true "Marker data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers/uuid/{uuid} [put]
func (h *MarkerHandler) UpdateMarkerByUUID(c *gin.Context) {
	ctx := c.Request.Context()

	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "UUID is required",
		})
		return
	}

	var marker models.Marker
	if err := c.ShouldBindJSON(&marker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	updatedMarker, err := h.markerRepo.UpdateMarkerByUUID(ctx, uuid, &marker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update marker",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Marker updated successfully",
		"data":    updatedMarker,
	})
}

// DeleteMarker godoc
// @Summary Delete marker by ID
// @Description Delete a marker by ID
// @Tags markers
// @Accept json
// @Produce json
// @Param id path int true "Marker ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers/{id} [delete]
func (h *MarkerHandler) DeleteMarker(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid marker ID",
			"error":   err.Error(),
		})
		return
	}

	err = h.markerRepo.DeleteMarker(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete marker",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Marker deleted successfully",
	})
}

// DeleteMarkerByUUID godoc
// @Summary Delete marker by UUID
// @Description Delete a marker by UUID
// @Tags markers
// @Accept json
// @Produce json
// @Param uuid path string true "Marker UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers/uuid/{uuid} [delete]
func (h *MarkerHandler) DeleteMarkerByUUID(c *gin.Context) {
	ctx := c.Request.Context()

	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "UUID is required",
		})
		return
	}

	err := h.markerRepo.DeleteMarkerByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete marker",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Marker deleted successfully",
	})
}

// GetMarkersWithPagination godoc
// @Summary Get markers with pagination
// @Description Get markers with pagination support
// @Tags markers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param sort query string false "Sort field" default("id")
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/markers/paginated [get]
func (h *MarkerHandler) GetMarkersWithPagination(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "")
	pageStr := c.DefaultQuery("page", "1")
	sortField := c.DefaultQuery("sort_field", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	search := c.Query("search")

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid limit parameter",
			"error":   "Limit must be a positive integer",
		})
		return
	}

	// Parse offset - if not provided, calculate from page
	var offset int
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid offset parameter",
				"error":   "Offset must be a non-negative integer",
			})
			return
		}
	} else {
		// Calculate offset from page
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
	}

	// Validate sort order
	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Create pagination struct with pointers
	pagination := models.Pagination{
		Limit:     &limit,
		Offset:    &offset,
		SortField: &sortField,
		SortOrder: &sortOrder,
	}

	// Add search if provided
	if search != "" {
		pagination.Search = &search
	}

	// Call repository method
	result, err := h.markerRepo.PageMarker(ctx, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get paginated markers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MarkerHandler) GetNearestMarkers(c *gin.Context) {
	imei := c.Query("imei")
	if imei == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IMEI is required"})
		return
	}

	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 10 {
		limit = 10
	}

	markers, err := h.markerRepo.GetNearestMarkers(c.Request.Context(), imei, lat, lng, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response dengan informasi tambahan
	response := gin.H{
		"data":      markers,
		"count":     len(markers),
		"timestamp": time.Now().Unix(),
		"coordinates": gin.H{
			"lat": lat,
			"lng": lng,
		},
		"imei": imei,
	}

	c.JSON(http.StatusOK, response)
}
