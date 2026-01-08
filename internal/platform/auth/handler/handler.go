package handler

import (
	db "bizbundl/internal/db/sqlc/platform"
	"bizbundl/internal/platform/auth/service"
	"bizbundl/internal/platform/auth/view"
	"bizbundl/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthHandler struct {
	service *service.AuthService
}

// CartService removed (Platform Users don't shop)
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Assuming UserResponse struct is defined elsewhere and has FirstName and LastName fields.
// Also assuming db.User struct has FirstName and LastName fields instead of FullName.
func newUserResponse(user db.User) UserResponse {
	// Platform User has no Phone
	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName, // Changed from FullName
		LastName:  user.LastName,  // Added
		// Phone:       "", // Not in schema
		Role: "owner", // Hardcode role
		// Permissions: nil,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}

	if req.Email == "" || req.Password == "" {
		return util.APIError(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid input"))
	}

	user, err := h.service.Register(c.Context(), req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// Auto-Login
	token, _, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		// Registration successful but login failed (rare/weird)
		// Return success but no token? Or error?
		// Let's return success but warn. Ideally shouldn't happen.
		return util.JSON(c, fiber.StatusCreated, user, "Registration successful (Login failed)")
	}

	h.setSessionCookie(c, token)

	return util.JSON(c, fiber.StatusCreated, fiber.Map{
		"token": token,
		"user":  newUserResponse(user),
	}, "Registration successful")
}

// ShowLoginForm renders the login page
func (h *AuthHandler) ShowLoginForm(c *fiber.Ctx) error {
	return util.Render(c, view.Login(view.LoginFormData{}))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse Form Data (HTMX sends x-www-form-urlencoded)
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		// Render form with error
		return util.Render(c, view.Login(view.LoginFormData{
			Message: "Invalid request format",
		}))
	}

	token, _, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		// Failure: Re-render form with error message
		return util.Render(c, view.Login(view.LoginFormData{
			Email:   req.Email,
			Message: "Invalid email or password", // Generic message
			// Errors: map[string]string{"email": "Invalid..."}, if we had field specific
		}))
	}

	h.setSessionCookie(c, token)

	// Success: HTMX Redirect
	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("session_token")
	return util.JSON(c, fiber.StatusOK, nil, "Logged out")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Middleware sets user_id from Token
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return util.APIError(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "Not authenticated"))
	}

	// Convert to UUID
	var uuid pgtype.UUID
	if err := uuid.Scan(userIDStr); err != nil {
		return util.APIError(c, fiber.StatusUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "Invalid User ID in token"))
	}

	// Fetch User from DB (DB Read happens ONLY here, not in middleware)
	// We need to expose a GetUser method in AuthService or use store directly if we had access?
	// Handler has access to Service. Let's add GetUser(ctx, id) to Service.

	user, err := h.service.GetUser(c.Context(), uuid)
	if err != nil {
		return util.APIError(c, fiber.StatusNotFound, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	return util.JSON(c, fiber.StatusOK, newUserResponse(user), "Current user")
}

func (h *AuthHandler) setSessionCookie(c *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * 7 * time.Hour)
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "Strict"
	c.Cookie(cookie)
}
