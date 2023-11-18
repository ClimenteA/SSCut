package routes

import (
	"github.com/ClimenteA/vidcastcutter/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func RoutesHandler(app *fiber.App) {
	r := app.Group("/")

	r.Post("/upload-video", handlers.UploadVideo)
	r.Post("/save-cutted", handlers.SaveCutted)
	r.Get("/metadata", handlers.GetVideoMetadata)
	r.Delete("/delete-video", handlers.DeleteVideo)
	r.Post("/export-video", handlers.ExportVideo)

}
