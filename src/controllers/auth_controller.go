package controllers

import (
	"api-gateway/src/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthServiceInterface
}

type AuthControllerInterface interface {
	RefreshToken(c *gin.Context)
	Login(c *gin.Context)
}
