package util

import (
	"math/rand"
	"os"
	"time"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
)

// DirFiles initializes the Doodads map with doodads.
func DirFiles(srcDir string) (filnames []string, err error) {
	folder, err := goutil.SrcDir(srcDir)
	if err != nil {
		return nil, errutil.Err(err)
	}
	f, err := os.Open(folder)
	if err != nil {
		return nil, errutil.Err(err)
	}
	fi, err := f.Readdir(0)
	if err != nil {
		return nil, errutil.Err(err)
	}

	var filenames []string
	for _, v := range fi {
		filename := folder + v.Name()
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

// RandInt is used by the debug function GenArea.
func RandInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
