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

func (h *Handler) PostRoomsFormulas(c echo.Context) error {
	var req models.PostRoomsFormulasJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// This is a simple mock that always returns a fixed board state.
	// It does not validate the formula or handle conflicts.
	mockBoard := models.Board{
		Content:   []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		Version:   2,
		GainScore: 10,
	}

	return c.JSON(http.StatusOK, mockBoard)
}

func (h *Handler) GetRoomsResult(c echo.Context) error {
	results := []models.RoomResultItem{
		{
			User:  "testuser1",
			Score: 120,
		},
		{
			User:  "testuser2",
			Score: 100,
		},
	}
	return c.JSON(http.StatusOK, results)
}
