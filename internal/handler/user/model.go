package user

type RegisterRequest struct {
	Username string `json:"username" example:"username"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword"`
}

type InfoResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
