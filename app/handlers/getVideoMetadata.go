package handlers

import (
	"log"

	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

func GetVideoPath() string {
	defaultVideoPath := "./videos/video.mp4"
	if funcs.PathExists(defaultVideoPath) && funcs.PathExists("./videos/frames") {
		return defaultVideoPath
	}
	return ""
}

type SkipFrames struct {
	Skip      int  `query:"skip"`
	FirstLoad bool `query:"firstLoad"`
}

type VidTSec struct {
	TotalSeconds int
}

func paginate(frames []types.FrameInfo, skip int) []types.FrameInfo {

	size := 7

	if skip < 0 {
		skip = 0
	}

	if skip > len(frames) {
		skip = len(frames)
	}

	end := skip + size
	if end > len(frames) {
		end = len(frames)
	}

	return frames[skip:end]

}

func GetVideoMetadata(c *fiber.Ctx) error {

	cfg, _ := c.Locals("cfg").(kvjson.DB)

	videoPath := GetVideoPath()
	if videoPath == "" {
		return c.Status(404).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "no video available",
		})
	}

	var totalSeconds int
	var tsec VidTSec
	cfg.Get("totalSeconds", &tsec)
	if tsec.TotalSeconds > 0 {
		totalSeconds = tsec.TotalSeconds
	} else {
		totalSeconds = funcs.GetVideoTotalSeconds(videoPath)
		cfg.SetWithKey("totalSeconds", VidTSec{TotalSeconds: totalSeconds})
	}

	kv, _ := c.Locals("kv").(kvjson.DB)
	frames := funcs.GetAllFramesSorted(kv)

	s := new(SkipFrames)
	if err := c.QueryParser(s); err != nil {
		return c.Status(400).JSON(types.WebResponse{
			Status:  string(types.Failed),
			Message: "skip not valid",
		})
	}

	if s.FirstLoad {
		cfg.Get("lastSkip", &s)
		frames = paginate(frames, s.Skip)
		log.Println("loaded frames from lastskip")
	} else {
		frames = paginate(frames, s.Skip)
		if s.Skip < 0 {
			s.Skip = 0
		}
		cfg.SetWithKey("lastSkip", s)
		log.Println("loaded fresh frames")
	}

	return c.Status(200).JSON(types.VideoMetadata{
		Status:       types.Success,
		Message:      "video available",
		TotalSeconds: totalSeconds,
		Frames:       frames,
		LastSkip:     s.Skip,
		VideoSrc:     "http://localhost:3000/videos/video.mp4",
	})

}
