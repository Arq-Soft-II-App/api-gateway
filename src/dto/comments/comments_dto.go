package comments

type CreateCommentDto struct {
	CourseId string `json:"course_id"`
	UserId   string `json:"user_id"`
	Text     string `json:"text"`
}

type GetCommentsDto struct {
	Text       string `json:"text"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	UserId     string `json:"user_id"`
}

type GetCommentsResponse []GetCommentsDto
