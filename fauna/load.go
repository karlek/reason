package fauna

import (
	"encoding/json"
	"io/ioutil"

	"github.com/karlek/reason/util"

	"github.com/karlek/worc/model"
	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

var Doodads = map[string]Doodad{}

// Load initializes the Doodads map with doodads.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/fauna/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		d, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		Doodads[d.Name()] = *d
	}
	return nil
}

// load parses a JSON data file into a go fauna object.
func load(filename string) (fa *Doodad, err error) {
	// jsonDoodad is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonDoodad struct {
		Name     string
		Graphics struct {
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

	var jc jsonDoodad
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

	fa = &Doodad{
		N: jc.Name,
		M: model.Model{
			G: termbox.Cell{
				Ch: rune(jc.Graphics.Ch[0]),
				Fg: fg,
				Bg: bg,
			},
			Pathable: jc.Pathable,
		},
	}
	return fa, nil
}
