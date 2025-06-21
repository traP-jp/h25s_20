package handler

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthUsecase    usecase.HealthUsecase
	roomUsecase      *usecase.RoomUsecase
	wsManager        *websocket.Manager
	WebSocketHandler *WebSocketHandler
}

func (h *Handler) GetHealth(c echo.Context) error {
	return h.HealthCheck(c)
}

func NewHandler(dbChecker domain.DatabaseHealthChecker, wsManager *websocket.Manager, roomUsecase *usecase.RoomUsecase) *Handler {
	wsHandler := NewWebSocketHandler(wsManager, roomUsecase)
	return &Handler{
		healthUsecase:    *usecase.NewHealthUsecase(dbChecker),
		roomUsecase:      roomUsecase,
		wsManager:        wsManager,
		WebSocketHandler: wsHandler,
	}
}
