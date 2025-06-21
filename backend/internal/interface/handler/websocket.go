package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WebSocketHandler struct {
	wsManager   domain.WebSocketManager
	roomManager domain.RoomManager
	wsUsecase   *usecase.WebSocketUsecase
}

func NewWebSocketHandler(wsManager domain.WebSocketManager, roomManager domain.RoomManager) *WebSocketHandler {
	wsUsecase := usecase.NewWebSocketUsecase(wsManager, roomManager)
	return &WebSocketHandler{
		wsManager:   wsManager,
		roomManager: roomManager,
		wsUsecase:   wsUsecase,
	}
}

// validateUserID validates and extracts user ID from request
func (h *WebSocketHandler) validateUserID(c echo.Context) (string, error) {
	// Try query parameter first
	userID := c.QueryParam("userId")
	if userID == "" {
		// Try header
		userID = c.Request().Header.Get("X-User-ID")
	}
	if userID == "" {
		// Try authorization header (for future JWT implementation)
		auth := c.Request().Header.Get("Authorization")
		if auth != "" {
			// TODO: Implement JWT token validation
			// For now, just log that we received auth header
			log.Debug().Str("auth", auth).Msg("Authorization header received but not yet implemented")
		}
	}

	if userID == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "User ID is required via 'userId' query parameter or 'X-User-ID' header")
	}

	// Basic validation
	if len(userID) < 1 || len(userID) > 100 {
		return "", echo.NewHTTPError(http.StatusBadRequest, "User ID must be between 1 and 100 characters")
	}

	return userID, nil
}

// HandleWebSocket handles WebSocket upgrade requests
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Validate and get user ID
	userID, err := h.validateUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("WebSocket upgrade failed: invalid user ID")
		return err
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := h.wsManager.UpgradeConnection(c.Response().Writer, c.Request(), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("Failed to upgrade WebSocket connection")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upgrade connection")
	}

	// Handle the connection in a separate goroutine with proper error recovery
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set a reasonable timeout for the connection
	go func() {
		timer := time.NewTimer(30 * time.Minute) // 30 minute timeout
		defer timer.Stop()

		select {
		case <-timer.C:
			log.Info().Str("userID", userID).Msg("WebSocket connection timeout")
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// Handle the WebSocket connection with panic recovery
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Str("userID", userID).Msg("Panic recovered in WebSocket handler")
			}
		}()

		h.wsUsecase.HandleConnection(ctx, conn)
	}()

	// Keep the handler alive until context is cancelled
	<-ctx.Done()
	return nil
}

// GetWebSocketStats returns WebSocket connection statistics (for debugging)
func (h *WebSocketHandler) GetWebSocketStats(c echo.Context) error {
	roomID := c.QueryParam("roomId")

	if roomID != "" {
		// Get stats for a specific room
		connections := h.wsManager.GetRoomConnections(roomID)
		users := make([]string, len(connections))
		for i, conn := range connections {
			users[i] = conn.GetUserID()
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"roomId":          roomID,
			"connectionCount": len(connections),
			"users":           users,
		})
	}

	// Return general stats (this is simplified - in a real app you'd want proper metrics)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "WebSocket service is running",
		"timestamp": time.Now().Unix(),
	})
}
