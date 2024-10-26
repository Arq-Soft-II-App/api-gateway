package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, controller controllers.AuthControllerInterface) {
	r.POST("/refresh-token", controller.RefreshToken)
	r.POST("/login", controller.Login)
}
