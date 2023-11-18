package handlers

import (
	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

func VideoTotalSeconds(c *fiber.Ctx) error {

	totalSeconds := funcs.GetVideoTotalSeconds("./videos/video.mp4")

	return c.JSON(types.TotalSecondsWebResponse{
		Status:       string(types.Success),
		TotalSeconds: totalSeconds,
	})

}
