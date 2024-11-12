package services

import "api-gateway/src/config/envs"

type Service struct {
	Auth       AuthServiceInterface
	Users      UsersServiceInterface
	Courses    CourseServiceInterface
	Comments   CommentsServiceInterface
	Categories CategoriesServiceInterface
	/* 	Inscriptions InscriptionsServiceInterface
	   	Search       SearchServiceInterface */
}

func NewService(env envs.Envs) *Service {
	usersService := NewUsersService(env)
	return &Service{
		Auth:       NewAuthService(env),
		Users:      usersService,
		Courses:    NewCourseService(env),
		Comments:   NewCommentsService(env, usersService),
		Categories: NewCategoriesService(env),
		/* 	Inscriptions: NewInscriptionsService(env),
		Search:       NewSearchService(env), */
	}
}
