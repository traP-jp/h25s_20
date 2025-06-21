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

func (h *Handler) PostRoomsActions(c echo.Context) error {
	var req models.PostRoomsActionsJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	switch req.Action {
	case models.JOIN, models.READY, models.CANCEL, models.START:
		// This is a simple mock that always returns success for valid actions.
		// It does not implement stateful logic for errors like 403 or 409.
		return c.NoContent(http.StatusNoContent)
	default:
		return c.JSON(http.StatusBadRequest, "Invalid action.")
	}
}
