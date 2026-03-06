package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSystemRoutes(router *gin.Engine) {

	router.GET("/systems", controllers.GetSystems)
}
