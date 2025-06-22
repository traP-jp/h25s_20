package main

import (
	"database/sql"
	_ "embed"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kaitoyama/kaitoyama-server-template/internal/db"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/auth"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/config"
	dbInfra "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/db"
	wsManager "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/interface/handler"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed openapi/swagger.yml
var swaggerFile []byte

func SetupRouter(database *sql.DB) *echo.Echo {
	// Load configuration for JWT secret
	cfg := config.LoadConfig()

	// Initialize Echo
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS設定を詳細に設定（WebSocket接続用）
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",              // フロントエンド開発サーバー
			"http://localhost:3000",              // 代替フロントエンドポート
			"http://app.frontend.orb.local:5173", // Docker環境用
		},
		AllowMethods: []string{
			echo.GET, echo.HEAD, echo.PUT, echo.PATCH,
			echo.POST, echo.DELETE, echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"Sec-WebSocket-Extensions",
			"Sec-WebSocket-Key",
			"Sec-WebSocket-Protocol",
			"Sec-WebSocket-Version",
			"Connection",
			"Upgrade",
		},
		AllowCredentials: true,
	}))

	// Load OpenAPI spec (use embedded file for buildpack compatibility)
	loader := openapi3.NewLoader()
	swagger, err := loader.LoadFromData(swaggerFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load swagger spec")
	}
	if err = swagger.Validate(loader.Context); err != nil {
		log.Fatal().Err(err).Msg("Swagger spec validation error")
	}

	// Initialize database queries
	queries := db.New(database)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userUsecase := usecase.NewUserUsecase(queries)
	authService := auth.NewAuthService(jwtService, userUsecase)
	roomUsecase := usecase.NewRoomUsecase()
	wsManagerInstance := wsManager.NewManager()
	wsHandler := handler.NewWebSocketHandler(wsManagerInstance, roomUsecase, userUsecase)

	// WebSocket endpoint (outside of API group to avoid OpenAPI validation)
	e.GET("/ws", wsHandler.HandleWebSocket)

	// Setup API routes
	api := e.Group("/api")

	dbChecker := dbInfra.NewDBHealthChecker(database)
	apiHandler := handler.NewHandler(dbChecker, wsManagerInstance, roomUsecase, userUsecase, jwtService)

	// 認証不要エンドポイント
	api.GET("/health", apiHandler.GetHealth)
	api.POST("/users", apiHandler.PostUsers)

	// 認証が必要なエンドポイント
	protectedApi := api.Group("")
	protectedApi.Use(authService.AuthMiddleware())

	protectedApi.GET("/rooms", apiHandler.GetRooms)
	protectedApi.POST("/rooms/:roomId/actions", func(c echo.Context) error {
		roomId, _ := strconv.Atoi(c.Param("roomId"))
		return apiHandler.PostRoomsRoomIdActions(c, roomId)
	})
	protectedApi.POST("/rooms/:roomId/formulas", func(c echo.Context) error {
		roomId, _ := strconv.Atoi(c.Param("roomId"))
		return apiHandler.PostRoomsRoomIdFormulas(c, roomId)
	})
	protectedApi.GET("/rooms/:roomId/result", func(c echo.Context) error {
		roomId, _ := strconv.Atoi(c.Param("roomId"))
		return apiHandler.GetRoomsRoomIdResult(c, roomId)
	})

	return e
}
