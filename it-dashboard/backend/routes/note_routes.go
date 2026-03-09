package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupNoteRoutes(router *gin.Engine) {

	//routes for notes
	router.GET("/api/systems/:id/notes", controllers.GetNotesBySystemID)
	router.POST("/api/systems/:id/notes", controllers.CreateNote)
}
