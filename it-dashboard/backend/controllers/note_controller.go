package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetNotesBySystemID godoc
// kemungkinan tidak diguankan karena sudah ada di system controller
// tapi tetap dibuat jika note banyak untuk bisa filter dan paginate
func GetNotesBySystemID(c *gin.Context) {
	systemID := c.Param("id")
	var notes []models.Note
	config.DB.Where("system_id = ?", systemID).Find(&notes)

	// Cek jika terjadi error pada database
	if config.DB.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil catatan"})
		return
	}
	// Return the notes as JSON response
	c.JSON(http.StatusOK, notes)
}

// CreateNote godoc
// @Summary Create a new note for a system
// @Description Create a new note for a system by providing the system ID and note content
func CreateNote(c *gin.Context) {
	systemIDStr := c.Param("id")
	// Konversi string ke uint64 dulu, lalu ke uint
	val, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID System tidak valid"})
		return
	}
	systemID := uint(val)

	// 1. CEK DULU: Apakah System dengan ID ini ada di database?
	var system models.System
	if err := config.DB.First(&system, systemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "System tidak ditemukan. Tidak bisa membuat catatan untuk sistem yang tidak ada.",
		})
		return
	}

	var request models.CreateNoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	note := models.Note{
		SystemId: systemID,
		Title:    request.Title,
		Content:  request.Content,
	}

	// Sebaiknya tangkap error jika database gagal simpan
	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, note)
}
