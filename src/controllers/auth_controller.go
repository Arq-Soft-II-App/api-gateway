package controllers

import (
	"api-gateway/src/dto/auth"
	"api-gateway/src/errors"
	"api-gateway/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthServiceInterface
}

func NewAuthController(service services.AuthServiceInterface) *AuthController {
	return &AuthController{
		service: service,
	}
}

type AuthControllerInterface interface {
	RefreshToken(c *gin.Context)
	Login(c *gin.Context)
}

func (c *AuthController) Login(ctx *gin.Context) {

	loginDTO := &auth.LoginDTO{}
	if err := ctx.BindJSON(loginDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar la solicitud",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	user, token, err := c.service.Login(*loginDTO)

	if err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			ctx.JSON(customErr.HTTPStatusCode, gin.H{
				"error": customErr.Message,
				"code":  customErr.Code,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}

	ctx.JSON(200, gin.H{"user": user, "token": token})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	user, newToken, err := c.service.RefreshToken(token)

	if err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			ctx.JSON(customErr.HTTPStatusCode, gin.H{
				"error": customErr.Message,
				"code":  customErr.Code,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}

	ctx.JSON(200, gin.H{"user": user, "token": newToken})
}
