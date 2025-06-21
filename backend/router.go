package main

import (
	"database/sql"
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/db"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/notification"
	wsManager "github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/websocket"
	"github.com/kaitoyama/kaitoyama-server-template/internal/interface/handler"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/kaitoyama/kaitoyama-server-template/openapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed openapi/swagger.yml
var swaggerFile []byte

// グローバルなNotificationServiceインスタンス
var NotificationService *notification.Service

func SetupRouter(database *sql.DB) *echo.Echo {
	// Initialize Echo
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Load OpenAPI spec (use embedded file for buildpack compatibility)
	loader := openapi3.NewLoader()
	swagger, err := loader.LoadFromData(swaggerFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load swagger spec")
	}
	if err = swagger.Validate(loader.Context); err != nil {
		log.Fatal().Err(err).Msg("Swagger spec validation error")
	}

	// Initialize services
	roomUsecase := usecase.NewRoomUsecase()
	wsManagerInstance := wsManager.NewManager()
	wsHandler := handler.NewWebSocketHandler(wsManagerInstance)

	// Initialize NotificationService (globally accessible)
	NotificationService = notification.NewService(roomUsecase, wsManagerInstance)

	// WebSocket endpoint (outside of API group to avoid OpenAPI validation)
	e.GET("/ws", wsHandler.HandleWebSocket)

	// Setup API routes
	api := e.Group("/api")

	dbChecker := db.NewDBHealthChecker(database)
	healthHandler := handler.NewHandler(dbChecker)

	// Register API handlers with /api prefix
	openapi.RegisterHandlersWithBaseURL(api, healthHandler, "")

	return e
}
