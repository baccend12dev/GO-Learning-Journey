package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDocumentationRoutes(router *gin.Engine) {
	// routes for documentation under a system
	documentationRoutes := router.Group("/api/systems/:id/documentations")
	{
		documentationRoutes.GET("", controllers.GetDocumentationsBySystemID)
		documentationRoutes.POST("", controllers.CreateDocumentation)
	}

	// routes for single documentation operations
	router.GET("/api/documentations/:id", controllers.GetDocumentationByID)
	router.PUT("/api/documentations/:id", controllers.UpdateDocumentation)
	router.DELETE("/api/documentations/:id", controllers.DeleteDocumentation)
}
