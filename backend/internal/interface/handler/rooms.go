package handler

import (
	"net/http"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRooms(c echo.Context) error {
	var res []models.Room
	rooms := h.roomUsecase.GetRooms()
	for _, room := range rooms {
		res = append(res, models.Room{
			RoomId:   room.RoomId,
			RoomName: room.RoomName,
			Users:    room.Users,
			IsOpened: room.IsOpened,
		})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostRoomsRoomIdActions(c echo.Context, roomId int) error {
	var req models.PostRoomsRoomIdActionsJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	var mockPlayer = domain.Player{
		ID:       1,
		UserName: "testuser",
	}

	switch req.Action {
	case models.JOIN:
		_, err := h.roomUsecase.AddPlayerToRoom(roomId, mockPlayer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to add player to room")
		}

		return c.NoContent(http.StatusNoContent)

	case models.READY:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, true)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to update player ready status")
		}
		return c.NoContent(http.StatusNoContent)
	case models.CANCEL:
		_, err := h.roomUsecase.UpdatePlayerReadyStatus(roomId, mockPlayer.ID, false)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to update player ready status")
		}
		return c.NoContent(http.StatusNoContent)
	case models.START:
		_, err := h.roomUsecase.StartGame(roomId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to start game")
		}
		return c.NoContent(http.StatusNoContent)
	case models.ABORT:
		_, err := h.roomUsecase.AbortGame(roomId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to abort game")
		}
		return c.NoContent(http.StatusNoContent)
	case models.CLOSERESULT:
		_, err := h.roomUsecase.CloseResult(roomId, mockPlayer.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to close result")
		}
		return c.NoContent(http.StatusNoContent)
	default:
		return c.JSON(http.StatusBadRequest, "Invalid action.")
	}
}

func (h *Handler) PostRoomsRoomIdFormulas(c echo.Context, roomId int) error {
	var req models.PostRoomsRoomIdFormulasJSONRequestBody
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

func (h *Handler) GetRoomsRoomIdResult(c echo.Context, roomId int) error {
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
