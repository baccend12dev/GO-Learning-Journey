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
