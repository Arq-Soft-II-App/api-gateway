package users

type UserResponseDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Lastname  string `json:"lastname"`
	Birthdate string `json:"birthdate"`
	Role      int    `json:"role"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}
