package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupServerRoutes(router *gin.Engine) {

	//group routes for server
	serverGroup := router.Group("/api/servers")
	{
		serverGroup.GET("/", controllers.GetServers)
		serverGroup.GET("/:id", controllers.GetServerByID)
		serverGroup.POST("/", controllers.CreateServer)
		serverGroup.PUT("/:id", controllers.UpdateServer)
		serverGroup.DELETE("/:id", controllers.DeleteServer)
	}
}
