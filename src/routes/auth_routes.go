package routes

import (
	"api-gateway/src/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(engine *gin.Engine, controller *controllers.AuthController) {
	engine.POST("/auth/refresh-token", controller.RefreshToken)
	engine.POST("/auth/login", controller.Login)
}
