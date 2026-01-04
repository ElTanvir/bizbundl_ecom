package auth

import (
	"bizbundl/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterRoutes sets up the API routes for Auth
func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	authGroup := router.Group("/auth")
	authGroup.Post("/register", h.Register)
	authGroup.Post("/login", h.Login)
	authGroup.Post("/logout", h.Logout)
	// authGroup.Get("/me", h.Me) // Middleware needed for this
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}

	// Note: Add validation here or in service. simple check for now
	if req.Email == "" || req.Password == "" {
		return util.APIError(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid input"))
	}

	user, err := h.service.Register(c.Context(), req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	return util.JSON(c, fiber.StatusCreated, user, "Registration successful")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}

	token, user, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return util.APIError(c, fiber.StatusUnauthorized, err)
	}

	// Set Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour) // Match service expiry
	cookie.HTTPOnly = true
	cookie.Secure = true // Assumption: Production uses HTTPS
	cookie.SameSite = "Strict"
	c.Cookie(cookie)

	return util.JSON(c, fiber.StatusOK, fiber.Map{
		"token": token,
		"user":  user,
	}, "Login successful")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Cookies("session_token")
	if token != "" {
		// Best effort logout
		_ = h.service.Logout(c.Context(), token)
	}

	// Clear Cookie
	c.ClearCookie("session_token")

	return util.JSON(c, fiber.StatusOK, nil, "Logged out")
}
