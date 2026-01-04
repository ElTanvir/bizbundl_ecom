package handler

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Permissions []byte `json:"permissions"`
	// Omit PasswordHash!
}
