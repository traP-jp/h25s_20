package handler

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthUsecase usecase.HealthUsecase
	wsHandler     *WebSocketHandler
	roomManager   domain.RoomManager
}

func (h *Handler) GetHealth(c echo.Context) error {
	return h.HealthCheck(c)
}

func NewHandler(dbChecker domain.DatabaseHealthChecker, wsManager domain.WebSocketManager, roomManager domain.RoomManager) *Handler {
	return &Handler{
		healthUsecase: *usecase.NewHealthUsecase(dbChecker),
		wsHandler:     NewWebSocketHandler(wsManager, roomManager),
		roomManager:   roomManager,
	}
}

// WebSocket related methods
func (h *Handler) HandleWebSocket(c echo.Context) error {
	return h.wsHandler.HandleWebSocket(c)
}

func (h *Handler) GetWebSocketStats(c echo.Context) error {
	return h.wsHandler.GetWebSocketStats(c)
}
