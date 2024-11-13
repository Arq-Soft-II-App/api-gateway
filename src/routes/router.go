/* package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, controller *controllers.Controller) {

	engine.OPTIONS("/*cors", func(c *gin.Context) {
		c.Status(204)
	})

	authGroup := engine.Group("/auth")
	{
		AuthRoutes(authGroup, controller.Auth)
	}

	userGroup := engine.Group("/users")
	{
		UserRoutes(userGroup, controller.Users)
	}

	courseGroup := engine.Group("/courses")
	{
		CoursesRoutes(courseGroup, controller.Courses, controller.Search)
	}
	commentsGroup := engine.Group("/comment")
	{
		CommentsRoutes(commentsGroup, controller.Comments)
	}

	ratingsGroup := engine.Group("/rating")
	{
		RatingsRoutes(ratingsGroup, controller.Ratings)
	}

	engine.POST("/enroll", middlewares.AuthMiddleware(), controller.Inscriptions.CreateInscription)

	engine.GET("/myCourses/", middlewares.AuthMiddleware(), controller.Inscriptions.GetMyCourses)

	engine.GET("/studentsInThisCourse/:cid", middlewares.AdminAuthMiddleware(), controller.Inscriptions.GetMyStudents)

	engine.GET("/isEnrolled/:cid", middlewares.AuthMiddleware(), controller.Inscriptions.IsAlreadyEnrolled)

	engine.POST("/category/create", controller.Categories.CreateCategory)
	engine.GET("/categories", controller.Categories.GetCategories)
}
*/

package routes

import (
	"api-gateway/src/controllers"
	"api-gateway/src/middlewares"
	"fmt"

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
	engine.GET("/courses", testCall("courses"), controller.Search.SearchCourses)
	engine.POST("/courses/create", controller.Courses.CreateCourse)
	engine.PUT("/courses/update", middlewares.AuthMiddleware(), controller.Courses.UpdateCourse)
	engine.GET("/courses/:id", controller.Courses.GetCourseById)

	// Rutas de comentarios
	engine.POST("/comment", controller.Comments.CreateComment)
	engine.GET("/comment/:cid", controller.Comments.GetCourseComments)
	engine.PUT("/comment", controller.Comments.UpdateComment)

	// Rutas de calificaciones
	engine.GET("/rating", testCall("rating"), controller.Ratings.GetAllRatings)
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

func testCall(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("testCall called", path)
	}
}
