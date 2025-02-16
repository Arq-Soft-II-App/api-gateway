package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes configura los endpoints para la administración de instancias.
func SetupAdminRoutes(engine *gin.Engine, adminController *controllers.AdminController) {
	// Creamos un grupo para rutas de admin y aplicamos el middleware de admin (ya lo tenés en adminMiddleware.go)
	adminGroup := engine.Group("/admin")
	adminGroup.Use(middlewares.AdminAuthMiddleware())

	// Endpoint para listar todas las instancias (contenedores)
	adminGroup.GET("/instances", adminController.ListInstances)
	// Endpoint para crear e iniciar una nueva instancia
	adminGroup.POST("/instances", adminController.StartInstance)
	// Endpoint para detener una instancia
	adminGroup.DELETE("/instances/:id", adminController.StopInstance)
}
