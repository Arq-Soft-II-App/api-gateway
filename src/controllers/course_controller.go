package controllers

import (
	"api-gateway/src/dto/courses"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	service services.CourseServiceInterface
}

func NewCourseController(service services.CourseServiceInterface) *CourseController {
	return &CourseController{
		service: service,
	}
}

type CourseControllerInterface interface {
	CreateCourse(c *gin.Context)
	UpdateCourse(c *gin.Context)
	GetCourseById(c *gin.Context)
	DeleteCourse(c *gin.Context)
}

func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var courseData courses.CourseDTO
	if err := ctx.ShouldBindJSON(&courseData); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("Datos inválidos"))
		return
	}

	result, err := c.service.CreateCourse(courseData)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	var courseData courses.CourseDTO
	if err := ctx.ShouldBindJSON(&courseData); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("Datos inválidos"))
		return
	}

	courseData.ID = ctx.Param("cid")

	result, err := c.service.UpdateCourse(courseData)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *CourseController) GetCourseById(ctx *gin.Context) {
	courseId := ctx.Param("id")

	course, err := c.service.GetCourseById(courseId)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, course)
}

func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	courseId := ctx.Param("id")

	err := c.service.DeleteCourse(courseId)
	if err != nil {
		ctx.JSON(errors.GetStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Curso eliminado"})
}
