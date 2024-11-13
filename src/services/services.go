package services

import "api-gateway/src/config/envs"

type Service struct {
	Auth         AuthServiceInterface
	Users        UsersServiceInterface
	Courses      CourseServiceInterface
	Comments     CommentsServiceInterface
	Ratings      RatingsServiceInterface
	Categories   CategoriesServiceInterface
	Inscriptions InscriptionServiceInterface
	Search       SearchServiceInterface
}

func NewService(env envs.Envs) *Service {
	usersService := NewUsersService(env)
	courseService := NewCourseService(env)
	return &Service{
		Auth:         NewAuthService(env),
		Users:        usersService,
		Courses:      courseService,
		Comments:     NewCommentsService(env, usersService),
		Ratings:      NewRatingsService(env),
		Categories:   NewCategoriesService(env),
		Inscriptions: NewInscriptionsService(env, courseService, usersService),
		Search:       NewSearchService(env),
	}
}
