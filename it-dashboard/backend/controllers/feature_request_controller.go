package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFeatureRequestsBySystemID retrieves all feature requests for a specific system
func GetFeatureRequestsBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	var featureRequests []models.FeatureRequest

	config.DB.Where("system_id = ?", systemID).Find(&featureRequests)

	// Check if database error occurred
	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data feature request"})
		return
	}

	c.JSON(http.StatusOK, featureRequests)
}

// CreateFeatureRequest creates a new feature request for a system
func CreateFeatureRequest(c *gin.Context) {
	systemIDStr := c.Param("id")
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. Check if System exists
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat request fitur untuk sistem yang tidak ada.",
		})
		return
	}

	// 2. Bind request body
	var request models.CreateFeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	featureRequest := models.FeatureRequest{
		SystemId:    systemID,
		Title:       request.Title,
		Description: request.Description,
		Status:      "Pending", // Default status is Pending
	}

	// 3. Save to database
	if err := config.DB.Create(&featureRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, featureRequest)
}

// UpdateFeatureRequest updates a feature request details (title, description, status)
func UpdateFeatureRequest(c *gin.Context) {
	var featureRequest models.FeatureRequest
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&featureRequest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature Request tidak ditemukan"})
		return
	}

	// 2. Bind payload
	var request models.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Prepare partial updates
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.Status != "" {
		updates["status"] = request.Status
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&featureRequest).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
			return
		}
	}

	// Fetch updated record
	config.DB.First(&featureRequest, id)
	c.JSON(http.StatusOK, featureRequest)
}

// DeleteFeatureRequest deletes a feature request by ID
func DeleteFeatureRequest(c *gin.Context) {
	var featureRequest models.FeatureRequest
	id := c.Param("id")

	// 1. Check if exists
	if err := config.DB.First(&featureRequest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature Request tidak ditemukan"})
		return
	}

	// 2. Delete record
	if err := config.DB.Delete(&featureRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feature Request deleted successfully",
		"id":      id,
	})
}
