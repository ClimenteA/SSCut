package funcs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
)

func GetAllFramesSorted(kv kvjson.DB) []types.FrameInfo {

	var frame types.FrameInfo
	var frames []types.FrameInfo

	entries, _ := os.ReadDir(kv.Path())

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		key := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		err := kv.Get(key, &frame)
		if err != nil {
			continue
		}
		frames = append(frames, frame)
	}

	sort.Slice(frames, func(i, j int) bool {
		return frames[i].Frame < frames[j].Frame
	})

	return frames

}
