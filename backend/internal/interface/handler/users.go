package handler

import (
	"net/http"

	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) PostUsers(c echo.Context) error {
	var user models.UserCreate
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	// TODO: Usecase層を呼び出し、永続化処理を実装する
	return c.NoContent(http.StatusCreated)
}
