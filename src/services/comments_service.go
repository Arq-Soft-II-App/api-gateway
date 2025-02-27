package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/comments"
	usersDto "api-gateway/src/dto/users"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CommentsService struct {
	CommentsInterface CommentsServiceInterface
	env               envs.Envs
	usersService      UsersServiceInterface
}

type CommentsServiceInterface interface {
	CreateComment(data comments.CreateCommentDto) (comments.CreateCommentDto, error)
	GetCourseComments(courseId string) (comments.GetCommentsResponse, error)
	UpdateComment(data comments.CreateCommentDto) (comments.GetCommentsDto, error)
}

type CourseComment struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

type CourseCommentsResponse []CourseComment

func NewCommentsService(env envs.Envs, usersService UsersServiceInterface) *CommentsService {
	return &CommentsService{
		env:          env,
		usersService: usersService,
	}
}

func (s *CommentsService) GetCourseComments(courseId string) (comments.GetCommentsResponse, error) {
	commentsURL := fmt.Sprintf("%s%s/%s", s.env.Get("COURSES_API_URL"), "ratings", courseId)

	req, err := http.NewRequest("GET", commentsURL, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud HTTP", err)
		return nil, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud HTTP", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al obtener los comentarios servicio", err)
		error := err.(*errors.Error)
		if error.HTTPStatusCode == 404 {
			return nil, errors.NewError("COMMENT_NOT_FOUND", "No se encontraron comentarios para el curso", 404)
		}
		fmt.Println("Error al obtener los comentarios", err)
		return nil, errors.NewError("REQUEST_ERROR", "Error al obtener los comentarios", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errors.NewError("COMMENT_NOT_FOUND", "No se encontraron comentarios para el curso", 404)
	}

	var courseComments CourseCommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&courseComments); err != nil {
		fmt.Println("Error al decodificar la respuesta comentarios", err)
		return nil, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta comentarios", http.StatusInternalServerError)
	}

	// Extraer los IDs de usuario únicos
	userIds := make([]string, 0)
	userIdsMap := make(map[string]bool)
	for _, comment := range courseComments {
		if !userIdsMap[comment.UserId] {
			userIds = append(userIds, comment.UserId)
			userIdsMap[comment.UserId] = true
		}
	}

	// Obtener información de usuarios
	users, err := s.usersService.GetUsersList(userIds)
	if err != nil {
		fmt.Println("Error al obtener los usuarios", err)
		return nil, err
	}

	// Crear mapa de usuarios para búsqueda rápida
	usersMap := make(map[string]usersDto.UserDTO)
	for _, user := range users {
		usersMap[user.ID] = user
	}

	// Construir respuesta final
	response := make(comments.GetCommentsResponse, len(courseComments))
	for i, comment := range courseComments {
		user := usersMap[comment.UserId]
		response[i] = comments.GetCommentsDto{
			Text:       comment.Text,
			UserName:   user.Name + " " + user.Lastname,
			UserAvatar: user.Avatar,
			UserId:     user.ID,
		}
	}

	return response, nil
}

func (s *CommentsService) CreateComment(data comments.CreateCommentDto) (comments.CreateCommentDto, error) {
	commentsURL := fmt.Sprintf("%s%s/", s.env.Get("COURSES_API_URL"), "ratings")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return comments.CreateCommentDto{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar el comentario", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", commentsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return comments.CreateCommentDto{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return comments.CreateCommentDto{}, errors.NewError("REQUEST_ERROR", "Error al crear el comentario", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var createdComment comments.CreateCommentDto
	if err := json.NewDecoder(resp.Body).Decode(&createdComment); err != nil {
		return comments.CreateCommentDto{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return createdComment, nil
}

func (s *CommentsService) UpdateComment(data comments.CreateCommentDto) (comments.GetCommentsDto, error) {
	commentsURL := fmt.Sprintf("%s%s/", s.env.Get("COURSES_API_URL"), "ratings")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return comments.GetCommentsDto{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar el comentario", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("PUT", commentsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return comments.GetCommentsDto{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return comments.GetCommentsDto{}, errors.NewError("REQUEST_ERROR", "Error al actualizar el comentario", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var updatedComment comments.GetCommentsDto
	if err := json.NewDecoder(resp.Body).Decode(&updatedComment); err != nil {
		return comments.GetCommentsDto{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	// Obtener información del usuario
	user, err := s.usersService.GetUsersList([]string{data.UserId})
	if err != nil {
		return comments.GetCommentsDto{}, err
	}

	// Unir información del usuario a la respuesta
	updatedComment.UserName = user[0].Name + " " + user[0].Lastname
	updatedComment.UserAvatar = user[0].Avatar

	return updatedComment, nil
}
