package middlewares

import (
	"strings"

	"api-gateway/src/config/envs"
	utilsJWT "api-gateway/src/utils/jwt"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
)

// AdminAuthMiddleware verifies the JWT token and the user's role
func AdminAuthMiddleware() gin.HandlerFunc {
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

		if claims.Role != 1 {
			ErrorResponse(c, 403, "You don't have permission to access this resource")
			return
		}

		c.Set("UserID", claims.Id)
		c.Next()
	}
}
