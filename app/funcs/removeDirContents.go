package funcs

import (
	"os"
	"path"
)

func RemoveDirContents(dirPath string) {
	dir, _ := os.ReadDir(dirPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{dirPath, d.Name()}...))
	}
}
