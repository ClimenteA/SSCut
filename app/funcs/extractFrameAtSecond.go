package funcs

import (
	"log"
	"os/exec"
	"strconv"
)

func ExtractFrameAtSecond(videoPath string, second int) bool {
	secondStr := strconv.Itoa(second)
	outputPath := "./videos/frames/" + secondStr + ".png"
	if PathExists(outputPath) {
		log.Println("Frame already extracted:", outputPath)
		return false
	}
	frameTimestamp := SecondsToHMS(second)
	cmd := exec.Command("ffmpeg", "-ss", frameTimestamp, "-i", videoPath, "-frames:v", "1", outputPath)
	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
		return false
	}
	log.Println(cmd)
	return true
}
