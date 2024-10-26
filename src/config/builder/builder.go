package builder

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/controllers"
	"api-gateway/src/routes"
	"api-gateway/src/services"

	"github.com/gin-gonic/gin"
)

func Build(env envs.Envs) *gin.Engine {
	service := services.NewService(env)
	controller := controllers.NewController(service)
	engine := gin.Default()
	routes.SetupRoutes(engine, controller)

	return engine
}
