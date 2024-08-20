package router

import (
	"github.com/gofiber/fiber/v2"

	core_handler "github.com/shordem/api.thryvo/handler/core"
	"github.com/shordem/api.thryvo/lib/config"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/middleware"
	core_repository "github.com/shordem/api.thryvo/repository/core"
	user_repository "github.com/shordem/api.thryvo/repository/user"
	core_service "github.com/shordem/api.thryvo/service/core"
)

func InitializeCoreRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// config
	fileConfig := config.NewFileConfig(env)

	// repository
	fileRepository := core_repository.NewFileRepository(db)
	folderRepository := core_repository.NewFolderRepository(db)
	userRepository := user_repository.NewUserRepository(db)

	// service
	fileService := core_service.NewFileService(fileConfig, fileRepository, folderRepository, userRepository)
	folderService := core_service.NewFolderService(folderRepository, userRepository)

	// handler
	fileHandler := core_handler.NewFileHandler(fileService)
	folderHandler := core_handler.NewFolderHandler(folderService)

	// Middlewares
	authMiddleware := middleware.Protected()
	apiKeyMiddleware := middleware.RequireAPIKey(db)

	// Base routes
	fileRouter := router.Group("/file")
	folderRouter := router.Group("/folder")

	fileRouter.Post("/upload", apiKeyMiddleware, fileHandler.UploadFile)
	fileRouter.Get("/", authMiddleware, fileHandler.GetUserFiles)
	fileRouter.Get("/:user_id/:key", fileHandler.GetFile)

	folderRouter.Post("/", authMiddleware, folderHandler.CreateFolder)
	folderRouter.Get("/", authMiddleware, folderHandler.GetUserFolders)
	folderRouter.Get("/:parent_id", authMiddleware, folderHandler.GetFoldersByParent)
	folderRouter.Put("/:id", authMiddleware, folderHandler.UpdateFolder)
	folderRouter.Delete("/:id", authMiddleware, folderHandler.DeleteFolder)
}
