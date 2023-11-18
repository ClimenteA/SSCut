package main

import (
	"github.com/ClimenteA/fiberwebgui"
	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const DEBUG = false

func main() {

	kv := kvjson.New("./.kvdb")
	cfgkv := kvjson.New("./.kvcfg")

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024 * 1024, // 100 GB
	})
	app.Static("/videos", "./videos", fiber.Static{
		ByteRange: true,
	})
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("kv", kv)
		c.Locals("cfg", cfgkv)
		return c.Next()
	})
	routes.RoutesHandler(app)

	defaultVideoPath := "./videos/video.mp4"
	if funcs.PathExists(defaultVideoPath) && funcs.PathExists("./videos/frames") {
		go funcs.ExtractFPS(defaultVideoPath, kv)
	}

	if DEBUG {
		app.Listen(":3000")
	} else {
		app.Static("/", "./ui/dist")
		fiberwebgui.RunBrowserOnPort(app, 3000)
	}

}
