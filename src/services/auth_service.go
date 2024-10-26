package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/auth"
	"api-gateway/src/dto/users"
	"api-gateway/src/errors"
	"api-gateway/src/utils/jwt"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type AuthService struct {
	AuthInterface AuthServiceInterface
	env           envs.Envs
}

type AuthServiceInterface interface {
	Login(data auth.LoginDTO) (users.UserDTO, string, error)
	RefreshToken(token string) (users.UserDTO, string, error)
}

func NewAuthService(env envs.Envs) *AuthService {
	return &AuthService{
		env: env,
	}
}

func (s *AuthService) Login(data auth.LoginDTO) (users.UserDTO, string, error) {
	Login_users_API_URL := s.env.Get("USERS_API_URL") + "/login"
	var userResponse users.UserDTO

	loginReq, err := json.Marshal(data)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("PAYLOAD_ERROR", "Error al serializar el payload", 500)
	}

	req, err := http.NewRequest("POST", Login_users_API_URL, bytes.NewBuffer(loginReq))
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("HTTP_REQUEST_ERROR", "Error al crear la solicitud HTTP", 500)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.env.Get("USERS_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("USERS_SERVICE_ERROR", "Error al realizar la solicitud al servicio de usuarios", 500)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("RESPONSE_READ_ERROR", "Error al leer la respuesta", 500)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return users.UserDTO{}, "", errors.NewError("INVALID_CREDENTIALS", "Credenciales inválidas", resp.StatusCode)
		}
		return users.UserDTO{}, "", errors.NewError("INVALID_CREDENTIALS", "Error al realizar la solicitud al servicio de usuarios", resp.StatusCode)
	}

	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		fmt.Println("Error al deserializar la respuesta: ", err)
		return users.UserDTO{}, "", errors.NewError("DESERIALIZATION_ERROR", "Error al deserializar la respuesta", 500)
	}

	userUUID, err := uuid.Parse(userResponse.ID)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("UUID_CONVERSION_ERROR", "Error al convertir el ID del usuario en UUID", 500)
	}

	token := jwt.SignDocument(userUUID, userResponse.Role)

	return userResponse, token, nil
}

func (a *AuthService) RefreshToken(token string) (users.UserDTO, string, error) {
	claims, err := jwt.VerifyToken(token)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("INVALID TOKEN", "Invalid token", 401)
	}

	fmt.Println(claims)
	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("INVALID ID", "Invalid ID", 401)
	}

	role := claims["role"].(string)

	getUserURL := fmt.Sprintf("%s/users/%s", a.env.Get("USERS_API_URL"), id)

	req, err := http.NewRequest("GET", getUserURL, nil)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("HTTP_REQUEST_ERROR", "Error al crear la solicitud HTTP", 500)
	}

	req.Header.Set("Authorization", a.env.Get("USERS_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("USERS_SERVICE_ERROR", "Error al realizar la solicitud al servicio de usuarios", 500)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return users.UserDTO{}, "", errors.NewError("USER_NOT_FOUND", "Usuario no encontrado", resp.StatusCode)
	}

	var checkUser users.UserDTO
	err = json.NewDecoder(resp.Body).Decode(&checkUser)
	if err != nil {
		return users.UserDTO{}, "", errors.NewError("DESERIALIZATION_ERROR", "Error al deserializar la respuesta", 500)
	}

	newToken := jwt.SignDocument(id, role)

	return checkUser, newToken, nil
}
