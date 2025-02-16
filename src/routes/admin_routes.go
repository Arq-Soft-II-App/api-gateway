package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(engine *gin.Engine, adminController *controllers.AdminController) {
	// Creamos un grupo para rutas de admin y aplicamos el middleware de admin
	adminGroup := engine.Group("/admin")
	adminGroup.Use(middlewares.AdminAuthMiddleware())

	// Endpoint para listar todas las instancias (contenedores)
	adminGroup.GET("/instances", adminController.ListInstances)
	// Endpoint para crear e iniciar una nueva instancia
	adminGroup.POST("/instances", adminController.CreatetInstance)
	// Endpoint para levantar una instancia
	adminGroup.POST("/instances/:id/start", adminController.StartInstance)
	// Endpoint para detener una instancia
	adminGroup.DELETE("/instances/:id", adminController.StopInstance)
	// Endpoint para remover una instancia
	adminGroup.DELETE("/instances/:id/remove", adminController.RemoveInstance)
}
