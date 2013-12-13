// Package beastiary contains information about all creatures in reason.
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

// Creatures is a map where names of the creature is the key mapping to that
// creature object.
var Creatures = map[string]Creature{}

// Creature is an object with a name.
type Creature struct {
	O object.Object
	N string
}

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

// jsonCreature is a temporary struct for easier conversion between JSON and
// go structs.
type jsonCreature struct {
	Name     string
	Graphics struct {
		Ch string
		Fg map[string]string
		Bg map[string]string
	}
	Stackable bool
}

// loadCreature parses a JSON data file into a go creature object.
func loadCreature(filename string) (c *Creature, err error) {
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
		N: jc.Name,
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

// Name returns the name of the creature.
func (c *Creature) Name() string {
	return c.N
}

// NewX sets a new x value for the coordinate.
func (c *Creature) NewX(x int) {
	c.O.NewX(x)
}

// NewY sets a new y value for the coordinate.
func (c *Creature) NewY(y int) {
	c.O.NewY(y)
}

// IsStackable returns whether objects can be stacked ontop of this object.
func (c *Creature) IsStackable() bool {
	return c.O.IsStackable()
}

// Graphic returns the graphic data of this object.
func (c *Creature) Graphic() termbox.Cell {
	return c.O.Graphic()
}

// X returns the x value of the current coordinate.
func (c *Creature) X() int {
	return c.O.X()
}

// Y returns the y value of the current coordinate.
func (c *Creature) Y() int {
	return c.O.Y()
}
