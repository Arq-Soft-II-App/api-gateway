package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/users"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UsersService struct {
	env envs.Envs
}

func NewUsersService(env envs.Envs) *UsersService {
	return &UsersService{
		env: env,
	}
}

type UsersServiceInterface interface {
	Register(user users.UserDTO) (users.UserDTO, error)
	Update(user users.UserDTO) (users.UserDTO, error)
}

func (s *UsersService) Register(user users.UserDTO) (users.UserDTO, error) {
	registerURL := s.env.Get("USERS_API_URL")

	client := &http.Client{}

	jsonData, err := json.Marshal(user)
	if err != nil {
		return users.UserDTO{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar los datos del usuario", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return users.UserDTO{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud HTTP", http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.env.Get("API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return users.UserDTO{}, errors.NewError("REQUEST_ERROR", "Error al enviar la solicitud al servicio de usuarios", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return users.UserDTO{}, errors.NewError("RESPONSE_READ_ERROR", "Error al leer la respuesta del servidor", http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusCreated {
		return users.UserDTO{}, errors.NewError("SERVER_ERROR", "Error en la respuesta del servidor", resp.StatusCode)
	}

	var createdUser users.UserDTO
	err = json.Unmarshal(body, &createdUser)
	if err != nil {
		return users.UserDTO{}, errors.NewError("DESERIALIZATION_ERROR", "Error al deserializar la respuesta del servidor", http.StatusInternalServerError)
	}

	return createdUser, nil
}

func (s *UsersService) Update(user users.UserDTO) (users.UserDTO, error) {
	updateURL := s.env.Get("USERS_API_URL")

	client := &http.Client{}

	jsonData, err := json.Marshal(user)
	if err != nil {
		return users.UserDTO{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar los datos del usuario", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return users.UserDTO{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud HTTP", http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.env.Get("API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return users.UserDTO{}, errors.NewError("REQUEST_ERROR", "Error al enviar la solicitud al servicio de usuarios", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return users.UserDTO{}, errors.NewError("RESPONSE_READ_ERROR", "Error al leer la respuesta del servidor", http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		return users.UserDTO{}, errors.NewError("SERVER_ERROR", fmt.Sprintf("Error en la respuesta del servidor: %s", string(body)), resp.StatusCode)
	}

	var updatedUser users.UserDTO
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return users.UserDTO{}, errors.NewError("DESERIALIZATION_ERROR", "Error al deserializar la respuesta del servidor", http.StatusInternalServerError)
	}

	return updatedUser, nil
}
