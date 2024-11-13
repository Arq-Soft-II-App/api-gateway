package inscriptions

type EnrollRequestResponseDto struct {
	CourseId string `json:"course_id"`
	UserId   string `json:"user_id"`
}
type Student struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

type CourseIdString struct {
	CourseId string `json:"course_id"`
}
type MyCourse struct {
	Id          string `json:"course_id"`
	CourseName  string `json:"course_name"`
	CourseImage string `json:"course_image"`
}

type CourseListDto struct {
	Courses []Course `json:"courses"`
}

type Course struct {
	CourseId string `json:"course_id"`
}

type StudentsInCourse []Student
type MyCourses []MyCourse
