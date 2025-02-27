package middlewares

import (
	"strings"

	"api-gateway/src/config/envs"
	utilsJWT "api-gateway/src/utils/jwt"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
)

// AuthMiddleware verifies the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(c, 400, "Authorization header is required")
			return
		}

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			ErrorResponse(c, 400, "Invalid token format")
			return
		}

		tokenString := tokenParts[1]

		claims := &utilsJWT.CustomClaims{}
		envs := envs.LoadEnvs(".env")
		secret := []byte(envs.Get("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(c, 401, "Invalid or expired token")
			return
		}

		c.Set("UserID", claims.Id)
		c.Next()
	}
}

func ErrorResponse(c *gin.Context, status int, message string) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
