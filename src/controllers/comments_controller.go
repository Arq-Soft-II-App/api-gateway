package controllers

import (
	"api-gateway/src/dto/comments"
	"api-gateway/src/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentsController struct {
	service services.CommentsServiceInterface
}

func NewCommentsController(service services.CommentsServiceInterface) *CommentsController {
	return &CommentsController{
		service: service,
	}
}

type CommentsControllerInterface interface {
	CreateComment(c *gin.Context)
	GetCourseComments(c *gin.Context)
	UpdateComment(c *gin.Context)
}

func (c *CommentsController) CreateComment(ctx *gin.Context) {
	var commentData comments.CreateCommentDto
	if err := ctx.ShouldBindJSON(&commentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	createdComment, err := c.service.CreateComment(commentData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "La comentario se registro con exito",
		"response": createdComment,
	})
}

func (c *CommentsController) GetCourseComments(ctx *gin.Context) {
	courseId := ctx.Param("cid")
	if courseId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de curso requerido"})
		return
	}

	comments, err := c.service.GetCourseComments(courseId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (c *CommentsController) UpdateComment(ctx *gin.Context) {
	var commentData comments.CreateCommentDto
	if err := ctx.ShouldBindJSON(&commentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updatedComment, err := c.service.UpdateComment(commentData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "El comentario se actualizo con exito",
		"response": updatedComment,
	})
}
