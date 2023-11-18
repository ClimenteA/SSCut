package funcs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetLastFrameSaved() int {

	files, err := os.ReadDir("./videos/frames")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return 0
	}

	var maxNum int
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}
		numStr := strings.TrimSuffix(name, ".png")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		if num > maxNum {
			maxNum = num
		}
	}

	return maxNum
}
