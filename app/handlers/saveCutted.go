package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

func SaveCutted(c *fiber.Ctx) error {

	kv, _ := c.Locals("kv").(kvjson.DB)

	var cutted types.FrameInfo
	if err := c.BodyParser(&cutted); err != nil {
		log.Println(err)
		return c.Status(400).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "Invalid cuttted body.",
		})
	}

	pngAvailable := funcs.PathExists(fmt.Sprintf("./videos/frames/%v.png", cutted.Frame))

	if cutted.Keep {
		if pngAvailable {
			cutted.Url = fmt.Sprintf("http://localhost:3000/videos/frames/%v.png", cutted.Frame)
		} else {
			cutted.Url = "/not-ready.svg"
		}
	}

	if !cutted.Keep {
		cutted.Url = "/frame-deleted.svg"
	}

	if err := kv.SetWithKey(strconv.Itoa(cutted.Frame), cutted); err != nil {
		log.Println(err)
		return c.Status(400).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "Can't save cutted.",
		})
	}

	return c.Status(200).JSON(types.WebResponse{
		Status:  string(types.Success),
		Message: "Saved current cut.",
	})

}
