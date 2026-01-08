package handler

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"` // Keep for backward compat input? Or force split?
	// User Requirement said "CAPI does better with First/Last".
	// We should probably accept full_name and split it internally for now to avoid breaking Frontend immediately,
	// OR update frontend.
	// The previous service code split it. So let's keep Input as FullName for now (easier migration)
	// BUT Output as First/Last.
	// Wait, the handler code `req.FullName` was used.
	// And service `Register` accepts `fullName` string.
	// So Input remains `FullName`.
	Phone string `json:"phone"`
}

// Actually, let's keep Input as FullName in Request to match Service signature for now.
// Service splits it.
// Response splits it.

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Permissions []byte `json:"permissions"`
	// Omit PasswordHash!
}
