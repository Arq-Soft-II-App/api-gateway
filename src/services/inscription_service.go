package services

import (
	"api-gateway/src/config/envs"
	"api-gateway/src/dto/courses"
	"api-gateway/src/dto/inscriptions"
	"api-gateway/src/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InscriptionService struct {
	InscriptionInterface InscriptionServiceInterface
	env                  envs.Envs
	courseService        CourseServiceInterface
	usersService         UsersServiceInterface
}

type InscriptionServiceInterface interface {
	CreateInscription(data inscriptions.EnrollRequestResponseDto) (inscriptions.EnrollRequestResponseDto, error)
	GetMyCourses(userId string) ([]courses.CourseListDto, error)
	GetCourseStudents(courseId string) (inscriptions.StudentsInCourse, error)
	IsEnrolled(courseId string, userId string) (bool, error)
}

func NewInscriptionsService(env envs.Envs, courseService CourseServiceInterface, usersService UsersServiceInterface) *InscriptionService {
	return &InscriptionService{
		env:           env,
		courseService: courseService,
		usersService:  usersService,
	}
}

func (s *InscriptionService) CreateInscription(data inscriptions.EnrollRequestResponseDto) (inscriptions.EnrollRequestResponseDto, error) {
	// Creamos canales para comunicar las goroutines.
	countCh := make(chan int, 1)                // Canal para recibir la cantidad de inscriptos.
	courseCh := make(chan courses.CourseDTO, 1) // Canal para recibir la info del curso.
	errCh := make(chan error, 2)                // Canal para capturar errores, por si alguno se pone loquita.

	// Goroutine para obtener la cantidad de inscriptos.
	// Esta va a buscar la lista de estudiantes inscriptos y luego contarla.
	go func() {
		students, err := s.GetCourseStudents(data.CourseId)
		if err != nil {
			errCh <- err // Si algo falla, avisamos por el canal de errores.
			return
		}
		countCh <- len(students) // Mandamos la cantidad de estudiantes.
	}()

	// Goroutine para obtener la información del curso.
	// Acá consultamos al servicio de cursos.
	go func() {
		course, err := s.courseService.GetCourseById(data.CourseId)
		if err != nil {
			errCh <- err // Avisamos si se rompe algo.
			return
		}
		courseCh <- course // Mandamos la info del curso.
	}()

	// Ahora esperamos a que lleguen ambos mensajes.
	var inscriptionsCount int
	var courseInfo courses.CourseDTO
	received := 0
	for received < 2 {
		select {
		// Si llega la cantidad de inscriptos, la guardamos.
		case cnt := <-countCh:
			inscriptionsCount = cnt
			received++
		// Si llega la información del curso, la guardamos.
		case course := <-courseCh:
			courseInfo = course
			received++
		// Si llega un error desde cualquiera de las goroutines, paramos todo y devolvemos el error y a lpm.
		case err := <-errCh:
			return inscriptions.EnrollRequestResponseDto{}, err
		}
	}

	// Verificamos que el curso tenga capacidad libre.
	// Es como ver si queda lugar en el bondi: si no hay, avisamos que el curso está lleno.
	if courseInfo.BaseCourseDto.CourseCapacity <= inscriptionsCount {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("COURSE_FULL", "El curso está completo", http.StatusBadRequest)
	}

	// Si el que se está inscribiendo es el último que cabe (cantidad de inscriptos == capacidad - 1),
	// lanzamos otra goroutine para actualizar el estado del curso y marcarlo como lleno (state = false).
	if inscriptionsCount == courseInfo.BaseCourseDto.CourseCapacity-1 {
		go func() {
			falseState := false
			courseInfo.BaseCourseDto.CourseState = &falseState
			_, err := s.courseService.UpdateCourse(courseInfo)
			if err != nil {
				// Si hay un error al actualizar, lo imprimimos, pero la inscripción ya sigue su curso.
				fmt.Printf("Error actualizando el estado del curso: %v\n", err)
			}
		}()
	}

	// Preparamos y ejecutamos la inscripción si todo va ok.
	inscriptionsURL := fmt.Sprintf("%senroll", s.env.Get("INSCRIPTIONS_URL"))
	jsonData, err := json.Marshal(data)
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("SERIALIZATION_ERROR", "Error al serializar la inscripción", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", inscriptionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("REQUEST_ERROR", "Error al crear la inscripción", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var createdInscription inscriptions.EnrollRequestResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&createdInscription); err != nil {
		fmt.Printf("Error al decodificar la respuesta: %v\n", err)
		return inscriptions.EnrollRequestResponseDto{}, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	// Si todo salió bien, devolvemos la inscripción creada.
	return createdInscription, nil
}

func (s *InscriptionService) GetMyCourses(userId string) ([]courses.CourseListDto, error) {
	fmt.Println("userId", userId)
	inscriptionsURL := fmt.Sprintf("%smyCourses?userId=%s", s.env.Get("INSCRIPTIONS_URL"), userId)

	resp, err := http.Get(inscriptionsURL)
	fmt.Println("resp inscriptions", resp)
	if err != nil {
		return nil, errors.NewError("REQUEST_ERROR", "Error al obtener las inscripciones", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Manejar el caso de 404
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewError("NOT_FOUND", "No hay inscripciones para este usuario", http.StatusNotFound)
	}

	// Decodificar la respuesta
	var inscriptionsData []inscriptions.Course
	if err := json.NewDecoder(resp.Body).Decode(&inscriptionsData); err != nil {
		fmt.Printf("Error al decodificar la respuesta: %v\n", err)
		return nil, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	// Extraer los IDs de los cursos
	var courseIDs []string
	for _, inscription := range inscriptionsData {
		courseIDs = append(courseIDs, inscription.CourseId)
	}

	// Llamar al servicio de cursos para obtener la información detallada
	coursesList, err := s.courseService.GetCoursesList(courseIDs)
	if err != nil {
		fmt.Printf("Error al obtener los cursos: %v\n", err)
		return nil, errors.NewError("COURSE_SERVICE_ERROR", "Error al obtener los cursos", http.StatusInternalServerError)
	}

	return coursesList, nil
}

func (s *InscriptionService) GetCourseStudents(courseId string) (inscriptions.StudentsInCourse, error) {
	inscriptionsURL := fmt.Sprintf("%sstudentsInThisCourse/%s", s.env.Get("INSCRIPTIONS_URL"), courseId)

	req, err := http.NewRequest("GET", inscriptionsURL, nil)
	if err != nil {
		return nil, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud HTTP", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewError("REQUEST_ERROR", "Error al obtener los estudiantes", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var students inscriptions.StudentsInCourse
	if err := json.NewDecoder(resp.Body).Decode(&students); err != nil {
		return nil, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	userIds := make([]string, 0)
	for _, student := range students {
		userIds = append(userIds, student.UserId)
	}

	// Obtener información de usuarios
	users, err := s.usersService.GetUsersList(userIds)
	if err != nil {
		return nil, err
	}

	// Construir respuesta final
	response := make(inscriptions.StudentsInCourse, len(users))
	for i, user := range users {
		response[i] = inscriptions.Student{
			UserId:   user.ID,
			UserName: user.Name + " " + user.Lastname,
			Avatar:   user.Avatar,
		}
	}

	return response, nil
}

func (s *InscriptionService) IsEnrolled(courseId string, userId string) (bool, error) {
	inscriptionsURL := fmt.Sprintf("%s/isEnrolled/%s/%s", s.env.Get("INSCRIPTIONS_URL"), courseId, userId)

	req, err := http.NewRequest("GET", inscriptionsURL, nil)
	if err != nil {
		return false, errors.NewError("REQUEST_CREATION_ERROR", "Error al crear la solicitud", http.StatusInternalServerError)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.NewError("REQUEST_ERROR", "Error al obtener la inscripción", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var response struct {
		Enrolled bool `json:"enrolled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error al decodificar la respuesta:", err)
		return false, errors.NewError("DECODE_ERROR", "Error al decodificar la respuesta", http.StatusInternalServerError)
	}

	return response.Enrolled, nil
}
