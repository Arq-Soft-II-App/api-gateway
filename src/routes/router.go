package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, controller *controllers.Controller) {

	// Manejar las solicitudes OPTIONS para CORS
	engine.OPTIONS("/*cors", func(c *gin.Context) {
		c.Status(204)
	})

	// Rutas de autenticación
	engine.POST("/auth/refresh-token", controller.Auth.RefreshToken)
	engine.POST("/auth/login", controller.Auth.Login)

	// Rutas de usuarios
	engine.POST("/users/register", controller.Users.Register)
	engine.PUT("/users/update", middlewares.AuthMiddleware(), controller.Users.Update)

	// Rutas de cursos
	engine.GET("/courses", controller.Search.SearchCourses)
	engine.POST("/courses/create", controller.Courses.CreateCourse)
	engine.PUT("/courses/update/:cid", middlewares.AuthMiddleware(), controller.Courses.UpdateCourse)
	engine.GET("/courses/:id", controller.Courses.GetCourseById)

	// Rutas de comentarios
	engine.POST("/comment", controller.Comments.CreateComment)
	engine.GET("/comment/:cid", controller.Comments.GetCourseComments)
	engine.PUT("/comment", controller.Comments.UpdateComment)

	// Rutas de calificaciones
	engine.GET("/rating", controller.Ratings.GetAllRatings)
	engine.POST("/rating", controller.Ratings.NewRating)
	engine.PUT("/rating", controller.Ratings.UpdateRating)

	// Rutas de inscripciones
	engine.GET("/myCourses/", middlewares.AuthMiddleware(), controller.Inscriptions.GetMyCourses)
	engine.GET("/studentsInThisCourse/:cid", middlewares.AdminAuthMiddleware(), controller.Inscriptions.GetMyStudents)
	engine.GET("/isEnrolled/:cid", middlewares.AuthMiddleware(), controller.Inscriptions.IsAlreadyEnrolled)
	engine.POST("/enroll", middlewares.AuthMiddleware(), controller.Inscriptions.CreateInscription)

	// Rutas de categorías
	engine.GET("/categories", controller.Categories.GetCategories)
	engine.POST("/category/create", controller.Categories.CreateCategory)
}
