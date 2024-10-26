package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, controller controllers.UsersControllerInterface) {
	r.POST("/register", controller.Register)
	r.PUT("/update", middlewares.AuthMiddleware(), controller.Update)
}
