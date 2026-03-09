package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.System{})
	config.DB.AutoMigrate(&models.Server{})
	config.DB.AutoMigrate(&models.Note{})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "database connection successful",
		})
	})

	routes.SetupSystemRoutes(r)
	routes.SetupServerRoutes(r)
	routes.SetupNoteRoutes(r)

	r.Run(":8080")
}
