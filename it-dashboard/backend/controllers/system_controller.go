package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSystems godoc
// @Summary Get all systems
// @Description Get a list of all systems
// GetSystems retrieves all systems from the database and sends them as a JSON response.
// It queries the database for all system records and returns them to the client.
// If no systems are found, it returns an empty list.
// If a database error occurs, it returns an appropriate error response.
func GetSystems(c *gin.Context) {
	var systems []models.System // Assuming you have a System model defined in your models package
	config.DB.Find(&systems)

	c.JSON(http.StatusOK, systems)
}

func GetSystemByID(c *gin.Context) {
	var system models.System
	id := c.Param("id")
	if err := config.DB.Preload("Server").First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	c.JSON(http.StatusOK, system)
}

// // create system database entry simple input example
// func CreateSystem(c *gin.Context) {
// 	var system models.System

// 	//validasi request body
// 	if err := c.ShouldBindJSON(&system); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// simpan data ke database
// 	config.DB.Create(&system)
// 	c.JSON(http.StatusCreated, system)

// }

// recomended create system database entry with request body struct
func CreateSystem(c *gin.Context) {
	var request models.CreateSystemRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	//create system sesuai dengan request body struct
	system := models.System{
		Name:        request.Name,
		Type:        request.Type,
		Links:       request.Links,
		ServerId:    request.ServerId,
		Status:      request.Status,
		Description: request.Description,
	}

	// simpan data ke database
	if err := config.DB.Create(&system).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create system",
			"details": err.Error(),
		})
		return
	}
	// return response dengan data system yang baru dibuat
	c.JSON(http.StatusCreated, gin.H{
		"message": "System created successfully",
		"system":  system,
	})

}

// update function for system
func UpdateSystem(c *gin.Context) {
	var system models.System
	id := c.Param("id")
	if err := config.DB.First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	if err := c.ShouldBindJSON(&system); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&system)
	c.JSON(http.StatusOK, system)
}

// DeleteSystem is a handler function that deletes a system from the database based on the provided ID.
// It retrieves the system record using the ID from the URL parameters, and if found, deletes it from the database.
// If the system is not found, it returns a 404 Not Found response. If the deletion is successful, it returns a 204 No Content response.

func DeleteSystem(c *gin.Context) {
	var system models.System
	id := c.Param("id")
	if err := config.DB.First(&system, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	config.DB.Delete(&system)
	c.JSON(http.StatusNoContent, nil)
}
