package auth

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
