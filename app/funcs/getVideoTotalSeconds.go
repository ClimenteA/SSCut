package funcs

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetVideoTotalSeconds(videoPath string) int {

	if !PathExists(videoPath) {
		return 0
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
	}

	secondsFloatStr := strings.TrimSpace(string(output))
	secondsFloat64, err := strconv.ParseFloat(secondsFloatStr, 64)
	if err != nil {
		log.Panicln(err)
	}
	secondsInt := int(secondsFloat64)

	return secondsInt
}
