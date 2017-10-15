package util

import (
	"math/rand"
	"os"
	"time"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

// Quit quits the game.
func Quit() {
	/// Add on release.
	// status.Print("Do you want to quit the game? [y/N]")
	// wantToQuit := util.NoOrYes()
	// if !wantToQuit {
	// 	return
	// }
	termbox.Close()
	os.Exit(0)
}

// YesOrNo forces the player to answer either y or n.
// Esc is false and enter is true.
func YesOrNo() bool {
	return prompt(true)
}

// NoOrYes forces the player to answer either y or n.
// Enter is false and esc is false.
func NoOrYes() bool {
	return prompt(false)
}

// prompt the user with y/n. Defaults esc to false, enter is variable.
func prompt(enter bool) bool {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'n', 'N':
				return false
			case 'y', 'Y':
				return true
			}
			switch ev.Key {
			case termbox.KeyEsc:
				return false
			case termbox.KeyEnter:
				return enter
			}
		}
	}
}

// DirFiles returns all filenames in a folder.
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
		filename := folder + "/" + v.Name()
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

// RandInt is used by the debug function GenArea.
func RandInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// ParseColor takes a JSON map that describes the color of a item and
// returns a termbox attribute.
func ParseColor(colorSetting map[string]string) (attr termbox.Attribute, err error) {
	if colorSetting == nil {
		return 0, nil
	}
	v, ok := colorSetting["color"]
	if !ok {
		return 0, errutil.Newf("missing map key `color` in: %v", colorSetting)
	}
	switch v {
	case "black":
		attr += termbox.ColorBlack
	case "red":
		attr += termbox.ColorRed
	case "green":
		attr += termbox.ColorGreen
	case "yellow":
		attr += termbox.ColorYellow
	case "blue":
		attr += termbox.ColorBlue
	case "magenta":
		attr += termbox.ColorMagenta
	case "cyan":
		attr += termbox.ColorCyan
	case "white":
		attr += termbox.ColorWhite
	}
	v, ok = colorSetting["attr"]
	if !ok {
		return 0, errutil.Newf("missing map key `attr`")
	}
	switch v {
	case "bold":
		attr += termbox.AttrBold
	case "underline":
		attr += termbox.AttrUnderline
	case "reverse":
		attr += termbox.AttrReverse
	}
	return attr, nil
}
