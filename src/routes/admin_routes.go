package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(engine *gin.Engine, adminController *controllers.AdminController) {
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

	// Endpoint para obtener los logs de un servicio
	// Ejemplos:
	//GET /admin/logs?service=nginx&since=3600
	//GET /admin/logs?since=2025-02-18T15:00:00&until=2025-02-18T16:00:00
	adminGroup.GET("/logs", adminController.GetLogs)

	// Nuevo endpoint para obtener las estadísticas de recursos de los contenedores
	adminGroup.GET("/stats", adminController.GetStats)
}
