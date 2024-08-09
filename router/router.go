package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/shordem/api.thryvo/handler"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/database"
)

func InitializeRouter(router *fiber.App, dbConn database.DatabaseInterface, env constants.Env) {

	main := router.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	main.Get("/monitor", monitor.New(monitor.Config{Title: "Thryvo API Monitor"}))

	InitializeUserRouter(main, dbConn, env)
	InitializeCoreRouter(main, dbConn, env)

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	router.Get("/logs/:key", func(c *fiber.Ctx) error {
		return handler.GetLogs(c, dbConn)
	})
	router.Get("/", handler.Index)
	router.Get("*", handler.NotFound)

}
