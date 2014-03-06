package creature

import (
	"encoding/json"
	"io/ioutil"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/util"

	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

// Beastiary is a map where names of the creature is the key mapping to that
// creature object.
/// TOTO mapping to basic creature types.
var Beastiary = map[string]Creature{}

// The protagonist.
var Hero Creature

// Load initializes the Beastiary map with creatures.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/creature/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		c, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		Beastiary[c.Name()] = *c
	}
	return nil
}

// load parses a JSON data file into a go creature object.
func load(filename string) (c *Creature, err error) {
	// jsonCreature is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonCreature struct {
		Name     string
		Graphics struct {
			Ch string
			Fg map[string]string
			Bg map[string]string
		}
		Hp           int
		Strength     int
		Sight        int
		Speed        int
		Regeneration int
		Pathable     bool
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

	fg, err := util.ParseColor(jc.Graphics.Fg)
	if err != nil {
		return nil, errutil.Err(err)
	}
	bg, err := util.ParseColor(jc.Graphics.Bg)
	if err != nil {
		return nil, errutil.Err(err)
	}

	c = &Creature{
		name:         jc.Name,
		Hp:           jc.Hp,
		MaxHp:        jc.Hp,
		Strength:     jc.Strength,
		Speed:        jc.Speed,
		Sight:        jc.Sight,
		Regeneration: jc.Regeneration,
		Inventory:    make(Inventory, len(item.Positions)),
	}
	c.SetPathable(jc.Pathable)
	c.SetGraphics(termbox.Cell{
		Ch: rune(jc.Graphics.Ch[0]),
		Fg: fg,
		Bg: bg,
	})
	return c, nil
}
