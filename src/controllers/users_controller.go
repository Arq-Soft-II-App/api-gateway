package controllers

import (
	"api-gateway/src/dto/users"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service services.UsersServiceInterface
}

func NewUsersController(service services.UsersServiceInterface) *UsersController {
	return &UsersController{
		service: service,
	}
}

type UsersControllerInterface interface {
	Register(c *gin.Context)
	Update(c *gin.Context)
}

func (c *UsersController) Register(ctx *gin.Context) {
	registerDTO := &users.UserDTO{}
	if err := ctx.BindJSON(registerDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar la solicitud",
			"code":  "INVALID_REQUEST",
		})
	}

	user, err := c.service.Register(*registerDTO)
	if err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			ctx.JSON(customErr.HTTPStatusCode, gin.H{
				"error": customErr.Message,
				"code":  customErr.Code,
			})
		}
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UsersController) Update(ctx *gin.Context) {
	updateDTO := &users.UserDTO{}
	if err := ctx.BindJSON(updateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar la solicitud",
			"code":  "INVALID_REQUEST",
		})
	}

	user, err := c.service.Update(*updateDTO)
	if err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			ctx.JSON(customErr.HTTPStatusCode, gin.H{
				"error": customErr.Message,
				"code":  customErr.Code,
			})
		}
	}

	ctx.JSON(http.StatusOK, user)
}
