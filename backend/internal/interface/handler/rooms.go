package handler

import (
	"net/http"

	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRooms(c echo.Context) error {
	rooms := []models.Room{
		{
			RoomId:   1,
			RoomName: "Room 1",
			Users:    []string{"player1", "player2"},
			IsOpened: true,
		},
		{
			RoomId:   2,
			RoomName: "Room 2",
			Users:    []string{"player3", "player4"},
			IsOpened: false,
		},
	}
	return c.JSON(http.StatusOK, rooms)
}
