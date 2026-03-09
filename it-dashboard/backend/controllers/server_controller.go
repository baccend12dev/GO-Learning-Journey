package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetServers godoc
// @Summary Get all servers
func GetServers(c *gin.Context) {
	var servers []models.Server
	config.DB.Find(&servers)
	c.JSON(http.StatusOK, servers)
}

func GetServerByID(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}

func CreateServer(c *gin.Context) {
	var request models.CreateServerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	server := models.Server{
		Name:     request.Name,
		IP:       request.IP,
		Location: request.Location,
		OS:       request.OS,
	}
	config.DB.Create(&server)
	c.JSON(http.StatusCreated, server)
}

// update server
func UpdateServer(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&server)
	c.JSON(http.StatusOK, server)
}

// Deleted server function, because we don't want to delete server data, we can just update the status to "inactive" or something like that. If you want to implement delete server function, you can uncomment the code below.
func DeleteServer(c *gin.Context) {
	var server models.Server
	id := c.Param("id")
	if err := config.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Id not found"})
		return
	}
	config.DB.Delete(&server)
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})
}
