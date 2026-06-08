package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupNoteRoutes(router *gin.Engine) {

	noteRoutes := router.Group("/api/systems/:id/notes")
	{
		noteRoutes.GET("", controllers.GetNotesBySystemID)
		noteRoutes.POST("", controllers.CreateNote)
	}

	router.PUT("/api/notes/:id", controllers.UpdateNote)
	router.DELETE("/api/notes/:id", controllers.DeleteNote)
}
