package handler

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/notification"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthUsecase       usecase.HealthUsecase
	roomUsecase         *usecase.RoomUsecase
	notificationService *notification.Service
}

func (h *Handler) GetHealth(c echo.Context) error {
	return h.HealthCheck(c)
}

func NewHandler(dbChecker domain.DatabaseHealthChecker, notificationService *notification.Service) *Handler {
	return &Handler{
		healthUsecase:       *usecase.NewHealthUsecase(dbChecker),
		roomUsecase:         usecase.NewRoomUsecase(),
		notificationService: notificationService,
	}
}
