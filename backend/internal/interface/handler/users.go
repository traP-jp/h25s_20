package handler

import (
	"net/http"

	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/kaitoyama/kaitoyama-server-template/openapi/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) PostUsers(c echo.Context) error {
	var user models.UserCreate
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	// まず既存ユーザーを確認
	_, err := h.userUsecase.GetUserByUsername(c.Request().Context(), user.Username)
	if err == nil {
		// 既存ユーザーの場合：ログイン処理
		authenticatedUser, authErr := h.userUsecase.AuthenticateUser(c.Request().Context(), user.Username, user.Password)
		if authErr != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid credentials"})
		}

		// JWTトークンを生成
		token, tokenErr := h.jwtService.GenerateToken(authenticatedUser.ID, authenticatedUser.Username)
		if tokenErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate token"})
		}

		// レスポンスを作成
		response := models.AuthResponse{
			Token: token,
			User: struct {
				Id       int    `json:"id"`
				Username string `json:"username"`
			}{
				Id:       int(authenticatedUser.ID),
				Username: authenticatedUser.Username,
			},
		}

		return c.JSON(http.StatusOK, response)
	}

	// 新規ユーザーの場合：登録処理
	createReq := usecase.CreateUserRequest{
		Username: user.Username,
		Password: user.Password,
	}

	createResp, createErr := h.userUsecase.CreateUser(c.Request().Context(), createReq)
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create user"})
	}

	// JWTトークンを生成
	token, tokenErr := h.jwtService.GenerateToken(int32(createResp.UserID), createResp.Username)
	if tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate token"})
	}

	// レスポンスを作成
	response := models.AuthResponse{
		Token: token,
		User: struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
		}{
			Id:       int(createResp.UserID),
			Username: createResp.Username,
		},
	}

	return c.JSON(http.StatusCreated, response)
}
