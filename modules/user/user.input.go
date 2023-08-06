package user

type CreateUserInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
