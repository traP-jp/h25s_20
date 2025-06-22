package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthService struct {
	jwtService  *JWTService
	userUsecase *usecase.UserUsecase
}

func NewAuthService(jwtService *JWTService, userUsecase *usecase.UserUsecase) *AuthService {
	return &AuthService{
		jwtService:  jwtService,
		userUsecase: userUsecase,
	}
}

func (a *AuthService) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-Forwarded-Userヘッダーをチェック
			forwardedUser := c.Request().Header.Get("X-Forwarded-User")
			if forwardedUser != "" {
				// X-Forwarded-Userが存在する場合の処理
				user, err := a.handleForwardedUser(c, forwardedUser)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "failed to process forwarded user")
				}

				// ユーザー情報をコンテキストに保存
				claims := &Claims{
					UserID:   int32(user.UserID),
					Username: user.Username,
				}
				ctx := context.WithValue(c.Request().Context(), UserContextKey, claims)
				c.SetRequest(c.Request().WithContext(ctx))

				return next(c)
			}

			// JWT認証を試行
			return a.handleJWTAuth(c, next)
		}
	}
}

func (a *AuthService) handleForwardedUser(c echo.Context, username string) (*usecase.CreateUserResponse, error) {
	// まず既存ユーザーを検索
	existingUser, err := a.userUsecase.GetUserByUsername(c.Request().Context(), username)
	if err == nil {
		// ユーザーが存在する場合
		return &usecase.CreateUserResponse{
			UserID:   int64(existingUser.ID),
			Username: existingUser.Username,
		}, nil
	}

	// ユーザーが存在しない場合、自動作成
	createReq := usecase.CreateUserRequest{
		Username: username,
		Password: "", // X-Forwarded-Userの場合はパスワード不要
	}

	return a.userUsecase.CreateUserWithoutPassword(c.Request().Context(), createReq)
}

func (a *AuthService) handleJWTAuth(c echo.Context, next echo.HandlerFunc) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
	}

	// Bearer tokenの形式チェック
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
	}

	token := parts[1]
	claims, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	// ユーザー情報をコンテキストに保存
	ctx := context.WithValue(c.Request().Context(), UserContextKey, claims)
	c.SetRequest(c.Request().WithContext(ctx))

	return next(c)
}

// ユーザー情報を取得するヘルパー関数
func GetUserFromContext(c echo.Context) (*Claims, bool) {
	user, ok := c.Request().Context().Value(UserContextKey).(*Claims)
	return user, ok
}
