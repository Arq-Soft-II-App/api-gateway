package services

import "api-gateway/src/config/envs"

type Service struct {
	Auth  AuthServiceInterface
	Users UsersServiceInterface
	/* 	Courses      CoursesServiceInterface
	   	Inscriptions InscriptionsServiceInterface
	   	Search       SearchServiceInterface */
}

func NewService(env envs.Envs) *Service {
	return &Service{
		Auth:  NewAuthService(env),
		Users: NewUsersService(env),
		/* 		Courses:      NewCoursesService(env),
		   		Inscriptions: NewInscriptionsService(env),
		   		Search:       NewSearchService(env), */
	}
}
