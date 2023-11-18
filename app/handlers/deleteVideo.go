package handlers

import (
	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

func DeleteVideo(c *fiber.Ctx) error {

	kv, _ := c.Locals("kv").(kvjson.DB)
	cfg, _ := c.Locals("cfg").(kvjson.DB)

	funcs.RemoveDirContents("videos")
	funcs.RemoveDirContents(kv.Path())
	funcs.RemoveDirContents(cfg.Path())

	return c.Status(200).JSON(types.WebResponse{
		Status:  string(types.Success),
		Message: "video deleted",
	})

}
