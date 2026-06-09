package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupFeatureRequestRoutes(router *gin.Engine) {
	featureRequestRoutes := router.Group("/api/systems/:id/feature-requests")
	featureRequestRoutes.Use(middleware.AuthMiddleware())
	{
		featureRequestRoutes.GET("", controllers.GetFeatureRequestsBySystemID)

		writeGroup := featureRequestRoutes.Group("")
		writeGroup.Use(middleware.RoleMiddleware("Administrator", "Developer"))
		{
			writeGroup.POST("", controllers.CreateFeatureRequest)
		}
	}

	featureIndividual := router.Group("/api/feature-requests/:id")
	featureIndividual.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Administrator", "Developer"))
	{
		featureIndividual.PUT("", controllers.UpdateFeatureRequest)
		featureIndividual.DELETE("", controllers.DeleteFeatureRequest)
	}
}
