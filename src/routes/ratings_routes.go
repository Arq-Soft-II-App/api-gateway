package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func RatingsRoutes(r *gin.RouterGroup, controller controllers.RatingsControllerInterface) {
	r.POST("/", controller.NewRating)
	r.PUT("/", controller.UpdateRating)
	r.GET("/", controller.GetAllRatings)
}
