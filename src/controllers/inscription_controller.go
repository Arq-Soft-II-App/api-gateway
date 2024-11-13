package controllers

import (
	"api-gateway/src/dto/inscriptions"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InscriptionsController struct {
	service services.InscriptionServiceInterface
}

func NewInscriptionsController(service services.InscriptionServiceInterface) *InscriptionsController {
	return &InscriptionsController{
		service: service,
	}
}

type InscriptionControllerInterface interface {
	CreateInscription(c *gin.Context)
	GetMyCourses(c *gin.Context)
	GetMyStudents(c *gin.Context)
	IsAlreadyEnrolled(c *gin.Context)
}

// CreateInscription crea una nueva inscripción
func (c *InscriptionsController) CreateInscription(ctx *gin.Context) {
	userId := ctx.GetString("UserID")
	var request inscriptions.EnrollRequestResponseDto
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewError("BIND_ERROR", "Error al enlazar el JSON", http.StatusBadRequest))
		return
	}

	request.UserId = userId

	createdInscription, err := c.service.CreateInscription(request)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "El usuario se registró con éxito", "response": createdInscription})
}

func (c *InscriptionsController) GetMyCourses(ctx *gin.Context) {
	userId := ctx.GetString("UserID")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, errors.NewError("INVALID_REQUEST", "Falta el userId", http.StatusBadRequest))
		return
	}

	courses, err := c.service.GetMyCourses(userId)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

// GetMyStudents retrieves the students of a specific course
func (c *InscriptionsController) GetMyStudents(ctx *gin.Context) {
	courseId := ctx.Param("cid")
	if courseId == "" {
		ctx.JSON(http.StatusBadRequest, errors.NewError("INVALID_REQUEST", "Falta el courseId", http.StatusBadRequest))
		return
	}

	students, err := c.service.GetCourseStudents(courseId)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, students)
}

func (c *InscriptionsController) IsAlreadyEnrolled(ctx *gin.Context) {
	courseId := ctx.Param("cid")
	userId := ctx.GetString("UserID")

	isEnrolled, err := c.service.IsEnrolled(courseId, userId)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	if !isEnrolled {
		ctx.JSON(http.StatusBadRequest, errors.NewError("USER_ALREADY_ENROLLED", "El usuario ya está inscrito", http.StatusBadRequest))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "El usuario no está inscrito"})
}
