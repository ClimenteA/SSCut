package funcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ClimenteA/vidcastcutter/app/types"
)

// ffmpeg -i video.mp4 -ss 00:00:05 -t 00:00:01 -c copy video-cutted.mp4
// ffmpeg: This is the command to invoke FFmpeg.
// -i video.mp4: This option specifies the input file for FFmpeg. In this case, it's video.mp4.
// -ss 00:00:05: This option specifies the starting point for the cut, indicated by the timestamp 00:00:05. It tells FFmpeg to start the cut from 5 seconds into the input video.
// -t 00:00:01: This option specifies the duration of the cut, indicated by the timestamp 00:00:01. It tells FFmpeg to cut only 1 second of the video starting from the specified starting point.
// -c copy: This option specifies the codec to be used for the output file. In this case, it's copy, which means FFmpeg will copy the audio and video streams without re-encoding them. This is a faster process compared to re-encoding, but it only works if the output container format supports the codec used in the input file.
// video-cutted.mp4: This is the output file name. FFmpeg will save the cut portion of the video to this file. The output file name can be chosen according to your preference.
func GetKeptSeconds(frames []types.FrameInfo) []int {

	var ks []int
	ss := 0
	for idx, f := range frames {
		if f.Keep {
			if idx == 0 && f.Frame > 0 {
				ss = 0
			} else {
				ss = f.Frame
			}
			ks = append(ks, ss)
		}
	}

	return ks

}

func GetVideoSlices(frames []types.FrameInfo) []string {
	var sst []string

	rawSeconds := GetKeptSeconds(frames)

	grouped := make([][]int, 0)
	temp := make([]int, 0)
	for i := 0; i < len(rawSeconds); i++ {
		temp = append(temp, rawSeconds[i])
		if i == len(rawSeconds)-1 || rawSeconds[i]+1 != rawSeconds[i+1] {
			grouped = append(grouped, temp)
			temp = make([]int, 0)
		}
	}

	for idx, gr := range grouped {
		ss := gr[0]
		t := 1
		if idx == 0 {
			t = len(gr) - 1
		} else {
			t = len(gr)
		}
		ssStr := SecondsToHMS(ss)
		tStr := SecondsToHMS(t)
		slice := fmt.Sprintf("-ss %s -t %s", ssStr, tStr)
		sst = append(sst, slice)
	}

	return sst

}

func CreateFFmpegCommands(frames []types.FrameInfo, options types.VideoOutputOptions) []string {

	keepFrames := GetVideoSlices(frames)

	ffmpegCmds := []string{}
	vidCuts := []string{}
	for idx, frames := range keepFrames {
		vidCut := fmt.Sprintf("./videos/concat/%d-video-cutted.mp4", idx)
		cmdStr := fmt.Sprintf("ffmpeg -i ./videos/video.mp4 %s -c copy %s", frames, vidCut)
		ffmpegCmds = append(ffmpegCmds, cmdStr)
		vidAbsCut, _ := filepath.Abs(vidCut)
		vidCuts = append(vidCuts, "file "+vidAbsCut)
	}

	ffmpegConcatCmd := ""
	ffmpegOutputVolCmd := ""
	if len(vidCuts) > 0 {
		vids := strings.Join(vidCuts, "\n")
		vidAbsExported, _ := filepath.Abs("./videos/video-exported.mp4")
		vidCutsPath := "./videos/concat/vidcuts.txt"
		os.WriteFile(vidCutsPath, []byte(vids), 0644)
		ffmpegConcatCmd = fmt.Sprintf("ffmpeg -f concat -safe 0 -i %s -c:v libx264 -c:a copy %s", vidCutsPath, vidAbsExported)
		if options.OutputVolume > 0 {
			vidVolAbsExported, _ := filepath.Abs("./videos/video-exported-increased-volume.mp4")
			// ffmpeg -i input.mp4 -af "volume=2.0" -c:v copy -c:a aac -b:a 256k output_higher_quality.mp4
			// ffmpegOutputVolCmd = fmt.Sprintf("ffmpeg -i %s -af \"volume=%d\" -c:v copy %s", vidAbsExported, options.OutputVolume, vidVolAbsExported)
			ffmpegOutputVolCmd = fmt.Sprintf("ffmpeg -i %s -af \"volume=%ddB\" -c:v copy %s", vidAbsExported, options.OutputVolume, vidVolAbsExported)
		}
	}

	if ffmpegConcatCmd != "" {
		ffmpegCmds = append(ffmpegCmds, ffmpegConcatCmd)
	}

	if ffmpegOutputVolCmd != "" {
		ffmpegCmds = append(ffmpegCmds, ffmpegOutputVolCmd)
	}

	fmt.Println("Generated the following ffmpeg commands:", ffmpegCmds)

	return ffmpegCmds

}
