package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// VerifyTokenResponse represents the response for token verification
type VerifyTokenResponse struct {
	Valid    bool   `json:"valid"`
	UserID   int32  `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

// VerifyToken verifies the JWT token validity
func (h *Handler) VerifyToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusOK, VerifyTokenResponse{
			Valid:   false,
			Message: "missing authorization header",
		})
	}

	// Bearer tokenの形式チェック
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusOK, VerifyTokenResponse{
			Valid:   false,
			Message: "invalid authorization header format",
		})
	}

	token := parts[1]
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return c.JSON(http.StatusOK, VerifyTokenResponse{
			Valid:   false,
			Message: "invalid or expired token",
		})
	}

	return c.JSON(http.StatusOK, VerifyTokenResponse{
		Valid:    true,
		UserID:   claims.UserID,
		Username: claims.Username,
		Message:  "token is valid",
	})
}
