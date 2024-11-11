package users

type UserDTO struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name"`
	Lastname  string  `json:"lastname"`
	Birthdate string  `json:"birthdate"`
	Role      string  `json:"role"`
	Email     string  `json:"email"`
	Avatar    string  `json:"avatar"`
	Password  *string `json:"password,omitempty"`
}

type UsersListDTO struct {
	Users []UserDTO `json:"users"`
}
