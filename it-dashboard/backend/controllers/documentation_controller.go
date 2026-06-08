package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDocumentationsBySystemID retrieves all documentation entries for a specific system (with optional category filter)
func GetDocumentationsBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	categoryQuery := c.Query("category")
	var docs []models.Documentation

	query := config.DB.Where("system_id = ?", systemID)
	if categoryQuery != "" {
		query = query.Where("category = ?", categoryQuery)
	}

	query.Find(&docs)

	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dokumentasi"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// GetDocumentationByID retrieves a single documentation record by ID
func GetDocumentationByID(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// CreateDocumentation creates a new documentation entry for a system
func CreateDocumentation(c *gin.Context) {
	systemIDStr := c.Param("id")
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. Verify System exists
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat dokumentasi untuk sistem yang tidak ada.",
		})
		return
	}

	// 2. Bind JSON
	var request models.CreateDocumentationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 3. Validate Category
	if !models.IsValidDocumentationCategory(request.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori dokumentasi tidak valid. Kategori yang diperbolehkan adalah: 'Business Flow', 'Technical Flow', 'API Documentation', 'Database Documentation', 'Deployment Guide', 'User Manual'",
		})
		return
	}

	doc := models.Documentation{
		SystemId: systemID,
		Title:    request.Title,
		Category: request.Category,
		Content:  request.Content,
	}

	// 4. Save to database
	if err := config.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// UpdateDocumentation updates a documentation entry
func UpdateDocumentation(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	// 2. Bind payload
	var request models.UpdateDocumentationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Validate category if provided
	if request.Category != "" && !models.IsValidDocumentationCategory(request.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori dokumentasi tidak valid. Kategori yang diperbolehkan adalah: 'Business Flow', 'Technical Flow', 'API Documentation', 'Database Documentation', 'Deployment Guide', 'User Manual'",
		})
		return
	}

	// 4. Apply updates
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Category != "" {
		updates["category"] = request.Category
	}
	if request.Content != "" {
		updates["content"] = request.Content
	}

	if len(updates) > 0 {
		if err := config.DB.Model(&doc).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
			return
		}
	}

	// Fetch updated record
	config.DB.First(&doc, id)
	c.JSON(http.StatusOK, doc)
}

// DeleteDocumentation deletes a documentation entry
func DeleteDocumentation(c *gin.Context) {
	var doc models.Documentation
	id := c.Param("id")

	// 1. Verify existence
	if err := config.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumentasi tidak ditemukan"})
		return
	}

	// 2. Delete record
	if err := config.DB.Delete(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus dokumentasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dokumentasi deleted successfully",
		"id":      id,
	})
}
