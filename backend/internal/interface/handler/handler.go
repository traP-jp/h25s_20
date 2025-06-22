package handler

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/auth"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthUsecase    usecase.HealthUsecase
	roomUsecase      *usecase.RoomUsecase
	userUsecase      *usecase.UserUsecase
	jwtService       *auth.JWTService
	wsManager        *websocket.Manager
	WebSocketHandler *WebSocketHandler
}

func (h *Handler) GetHealth(c echo.Context) error {
	return h.HealthCheck(c)
}

func NewHandler(dbChecker domain.DatabaseHealthChecker, wsManager *websocket.Manager, roomUsecase *usecase.RoomUsecase, userUsecase *usecase.UserUsecase, jwtService *auth.JWTService) *Handler {
	wsHandler := NewWebSocketHandler(wsManager, roomUsecase)
	return &Handler{
		healthUsecase:    *usecase.NewHealthUsecase(dbChecker),
		roomUsecase:      roomUsecase,
		userUsecase:      userUsecase,
		jwtService:       jwtService,
		wsManager:        wsManager,
		WebSocketHandler: wsHandler,
	}
}
