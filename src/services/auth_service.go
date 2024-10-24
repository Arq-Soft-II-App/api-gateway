package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/auth"
	"api-gateway/src/dto/users"
	"api-gateway/src/errors"
	"api-gateway/src/utils/jwt"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type AuthService struct {
	AuthInterface AuthServiceInterface
	env           envs.Envs
}

type AuthServiceInterface interface {
	Login(email string, password string) (users.UserResponseDTO, string, error)
	RefreshToken(token string) (string, error)
}

func NewAuthService(env envs.Envs) *AuthService {
	return &AuthService{
		env: env,
	}
}

func (s *AuthService) Login(email string, password string) (users.UserResponseDTO, string, error) {
	usersAPIURL := s.env.Get("USERS_API_URL") + "/login"
	var userResponse users.UserResponseDTO

	loginReq, err := json.Marshal(auth.LoginDTO{Email: email, Password: password})
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("PAYLOAD_ERROR", "Error al serializar el payload", 500)
	}

	req, err := http.NewRequest("POST", usersAPIURL, bytes.NewBuffer(loginReq))
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("HTTP_REQUEST_ERROR", "Error al crear la solicitud HTTP", 500)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.env.Get("USERS_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("USERS_SERVICE_ERROR", "Error al realizar la solicitud al servicio de usuarios", 500)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("RESPONSE_READ_ERROR", "Error al leer la respuesta", 500)
	}

	if resp.StatusCode != http.StatusOK {

		return users.UserResponseDTO{}, "", errors.NewError("INVALID_CREDENTIALS", "Credenciales inválidas o error en el servicio de usuarios", resp.StatusCode)
	}

	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("DESERIALIZATION_ERROR", "Error al deserializar la respuesta", 500)
	}

	userUUID, err := uuid.Parse(userResponse.ID)
	if err != nil {
		return users.UserResponseDTO{}, "", errors.NewError("UUID_CONVERSION_ERROR", "Error al convertir el ID del usuario en UUID", 500)
	}

	token := jwt.SignDocument(userUUID, userResponse.Role)

	return userResponse, token, nil
}

func (s *AuthService) RefreshToken(token string) (string, error) {
	return "", nil
}
