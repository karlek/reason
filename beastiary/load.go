package beastiary

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/karlek/worc/object"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

// LoadCreatures initializes the Creatures map with creatures.
func LoadCreatures() (err error) {
	folder, err := goutil.SrcDir("github.com/karlek/reason/beastiary/data/")
	if err != nil {
		return errutil.Err(err)
	}
	f, err := os.Open(folder)
	if err != nil {
		return errutil.Err(err)
	}
	fi, err := f.Readdir(0)
	if err != nil {
		return errutil.Err(err)
	}

	for _, v := range fi {
		filename := folder + v.Name()
		c, err := loadCreature(filename)
		if err != nil {
			return errutil.Err(err)
		}
		Creatures[c.Name()] = *c
	}
	return nil
}

// loadCreature parses a JSON data file into a go creature object.
func loadCreature(filename string) (c *Creature, err error) {
	// jsonCreature is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonCreature struct {
		Name     string
		Graphics struct {
			Ch string
			Fg map[string]string
			Bg map[string]string
		}
		Hp        int
		Strength  int
		Speed     float64
		Stackable bool
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errutil.Err(err)
	}

	var jc jsonCreature
	err = json.Unmarshal(buf, &jc)
	if err != nil {
		return nil, errutil.Err(err)
	}

	fg, err := parseColor(jc.Graphics.Fg)
	if err != nil {
		return nil, errutil.Err(err)
	}
	bg, err := parseColor(jc.Graphics.Bg)
	if err != nil {
		return nil, errutil.Err(err)
	}

	c = &Creature{
		N:        jc.Name,
		Hp:       jc.Hp,
		MaxHp:    jc.Hp,
		Strength: jc.Strength,
		Speed:    jc.Speed,
		O: object.Object{
			G: termbox.Cell{
				Ch: rune(jc.Graphics.Ch[0]),
				Fg: fg,
				Bg: bg,
			},
			Stackable: jc.Stackable,
		},
	}
	return c, nil
}

// parseColor takes a JSON map that describes the color of a creature and
// returns a termbox attribute.
func parseColor(colorSetting map[string]string) (attr termbox.Attribute, err error) {
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
