package item

import (
	"encoding/json"
	"io/ioutil"

	"github.com/karlek/reason/util"

	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

// Items is a map where names of the item is the key mapping to that
// item object.
var Items = map[string]*Item{}

// Load initializes the Items map with creatures.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/item/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		i, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		Items[i.Name()] = i
	}
	return nil
}

// load parses a JSON data file into a go item object.
func load(filename string) (i *Item, err error) {
	// jsonItem is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonItem struct {
		Name        string
		Category    string
		Description string
		Num         int
		Graphics    struct {
			Ch string
			Fg map[string]string
			Bg map[string]string
		}
		Pathable bool
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errutil.Err(err)
	}

	var jc jsonItem
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

	j := new(Item)
	// j = &Item{
	// 	name:        jc.Name,
	// 	Category:    jc.Category,
	// 	Description: jc.Description,
	// 	Num:         jc.Num,
	// }

	j.name = jc.Name
	j.Category = jc.Category
	j.Description = jc.Description
	j.Num = jc.Num

	// if jc.Category == "weapon" {
	// 	i = &Weapon{}
	// }

	j.SetPathable(jc.Pathable)
	j.SetGraphics(termbox.Cell{
		Ch: rune(jc.Graphics.Ch[0]),
		Fg: fg,
		Bg: bg,
	})
	return j, nil
}
