package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func InscriptionRoutes(g *gin.RouterGroup, controller controllers.InscriptionControllerInterface) {
	g.POST("/enroll",
		middlewares.AuthMiddleware(),
		controller.CreateInscription)

	g.GET("/myCourses/",
		middlewares.AuthMiddleware(),
		controller.GetMyCourses)

	g.GET("/studentsInThisCourse/:cid",
		middlewares.AdminAuthMiddleware(),
		controller.GetMyStudents)

	g.GET("/isEnrolled/:cid",
		middlewares.AuthMiddleware(),
		controller.IsAlreadyEnrolled)
}
