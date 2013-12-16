package util

import (
	"os"

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
