package funcs

import (
	"testing"

	"github.com/ClimenteA/vidcastcutter/app/types"
)

func TestGenerateCMDs(t *testing.T) {
	options := types.VideoOutputOptions{}

	frames := []types.FrameInfo{
		{Frame: 0, Keep: true, Url: ""},
		{Frame: 1, Keep: true, Url: ""},
		{Frame: 2, Keep: true, Url: ""},
		{Frame: 3, Keep: false, Url: ""},
		{Frame: 4, Keep: false, Url: ""},
		{Frame: 5, Keep: false, Url: ""},
		{Frame: 6, Keep: true, Url: ""},
		{Frame: 7, Keep: true, Url: ""},
		{Frame: 8, Keep: false, Url: ""},
		{Frame: 9, Keep: true, Url: ""},
		{Frame: 10, Keep: true, Url: ""},
		{Frame: 11, Keep: true, Url: ""},
		{Frame: 12, Keep: false, Url: ""},
		{Frame: 13, Keep: true, Url: ""},
		{Frame: 14, Keep: true, Url: ""},
		{Frame: 15, Keep: true, Url: ""},
	}

	expected := []string{
		"-ss 00:00:00 -t 00:00:02",
		"-ss 00:00:06 -t 00:00:02",
		"-ss 00:00:09 -t 00:00:03",
		"-ss 00:00:13 -t 00:00:03",
	}

	result := CreateFFmpegCommands(frames, options)

	t.Errorf("Result: %s,\n Expected: %s", result, expected)
}
