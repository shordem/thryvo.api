package router

import (
	"github.com/gofiber/fiber/v2"

	core_handler "github.com/shordem/api.thryvo/handler/core"
	"github.com/shordem/api.thryvo/lib/config"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/middleware"
)

func InitializeCoreRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// service
	mediaService := config.NewFileHelper(env)

	// handler
	mediaHandler := core_handler.NewMediaHandler(mediaService)

	// Middlewares
	authMiddleware := middleware.Protected()

	// Base routes
	mediaRouter := router.Group("/media")

	mediaRouter.Post("/upload", mediaHandler.UploadMedia, authMiddleware)
	mediaRouter.Get("/:key", mediaHandler.GetMedia)
}
