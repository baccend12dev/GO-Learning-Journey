package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupFeatureRequestRoutes(router *gin.Engine) {
	featureRequestRoutes := router.Group("/api/systems/:id/feature-requests")
	{
		featureRequestRoutes.GET("", controllers.GetFeatureRequestsBySystemID)
		featureRequestRoutes.POST("", controllers.CreateFeatureRequest)
	}

	router.PUT("/api/feature-requests/:id", controllers.UpdateFeatureRequest)
	router.DELETE("/api/feature-requests/:id", controllers.DeleteFeatureRequest)
}
