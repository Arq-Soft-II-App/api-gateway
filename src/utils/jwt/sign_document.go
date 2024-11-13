package jwt

import (
	config "api-gateway/src/config/envs"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Id   uuid.UUID `json:"id"`
	Role int       `json:"role"`
}

func (c *CustomClaims) Valid() error {
	return nil
}

func NewCustomClaims(id uuid.UUID, role string) *CustomClaims {
	return &CustomClaims{
		Id:   id,
		Role: setRole(role),
	}
}

func SignDocument(id uuid.UUID, role string) string {
	claims := NewCustomClaims(id, role)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	envs := config.LoadEnvs(".env")
	secret := []byte(envs.Get("JWT_SECRET"))
	signedToken, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return signedToken
}

func setRole(role string) int {
	if role == "admin" {
		return 1
	}
	return 0
}
