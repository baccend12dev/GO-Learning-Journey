package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSystemRoutes(router *gin.Engine) {

	//group routes for system
	systemGroup := router.Group("/api/systems")
	{
		systemGroup.GET("/", controllers.GetSystems)
		systemGroup.GET("/:id", controllers.GetSystemByID)
		systemGroup.POST("/", controllers.CreateSystem)
	}

}
