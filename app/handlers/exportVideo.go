package handlers

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

func ExportVideo(c *fiber.Ctx) error {

	kv, _ := c.Locals("kv").(kvjson.DB)

	exportedPath := "./videos/video-exported.mp4"
	if funcs.PathExists(exportedPath) {
		os.Remove(exportedPath)
	}
	exportedVidVolPath := "./videos/video-exported-increased-volume.mp4"
	if funcs.PathExists(exportedVidVolPath) {
		os.Remove(exportedVidVolPath)
	}

	frames := funcs.GetAllFramesSorted(kv)

	var options types.VideoOutputOptions
	if err := c.BodyParser(&options); err != nil {
		log.Println(err)
		return c.Status(400).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "Invalid cuttted body.",
		})
	}

	ffmpegCommands := funcs.CreateFFmpegCommands(frames, options)

	log.Println("FFmpeg Commands:", ffmpegCommands)

	if len(ffmpegCommands) == 0 {
		return c.Status(400).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "nothing to cut",
		})
	}

	for _, command := range ffmpegCommands {
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to execute command: %s\nError: %s\n", command, err)
		} else {
			fmt.Printf("Command executed successfully: %s\nOutput: %s\n", command, string(output))
		}
	}

	return c.Status(200).JSON(types.WebResponse{
		Status:  string(types.Success),
		Message: "video exported",
	})

}
