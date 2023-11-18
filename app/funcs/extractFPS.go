package funcs

import (
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
)

func ExtractFPS(videoPath string, kv kvjson.DB) int {
	totalSeconds := GetVideoTotalSeconds(videoPath)
	log.Printf("Video has %d seconds. Started extracting FPS...", totalSeconds)

	for second := 1; second <= totalSeconds; second++ {
		cutted := types.FrameInfo{
			Frame: second,
			Keep:  true,
			Url:   "/not-ready.svg",
		}
		kv.SetWithKeyIfNew(strconv.Itoa(cutted.Frame), cutted)
	}

	cores := runtime.NumCPU()
	if cores >= 4 {
		cores = 2
	} else {
		cores = 1
	}
	limiter := make(chan int, cores)
	for second := 1; second <= totalSeconds; second++ {
		limiter <- 1
		go func(second int) {
			extracted := ExtractFrameAtSecond(videoPath, second)
			if extracted {
				var existingCutted types.FrameInfo
				kv.Get(strconv.Itoa(second), &existingCutted)
				url := fmt.Sprintf("http://localhost:3000/videos/frames/%v.png", second)
				if !existingCutted.Keep {
					url = "/frame-deleted.svg"
				}
				cutted := types.FrameInfo{
					Frame: second,
					Keep:  existingCutted.Keep,
					Url:   url,
				}
				kv.SetWithKey(strconv.Itoa(cutted.Frame), cutted)
			}
			<-limiter
		}(second)
	}

	log.Println("Done extracting FPS!")
	return 1
}
