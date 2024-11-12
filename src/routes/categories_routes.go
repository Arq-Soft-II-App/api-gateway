package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func CategoriesRoutes(r *gin.RouterGroup, controller controllers.CategoriesControllerInterface) {
	r.POST("/", controller.CreateCategory)
	r.GET("/", controller.GetCategories)
}
