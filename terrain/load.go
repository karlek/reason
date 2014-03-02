package terrain

import (
	"encoding/json"
	"io/ioutil"

	"github.com/karlek/reason/util"

	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

var Fauna = map[string]Terrain{}

// Load initializes the fauna collection.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/terrain/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		t, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		Fauna[t.Name()] = *t
	}
	return nil
}

// load parses a JSON data file into a go terrain object.
func load(filename string) (t *Terrain, err error) {
	// jsonTerrain is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonTerrain struct {
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

	var jt jsonTerrain
	err = json.Unmarshal(buf, &jt)
	if err != nil {
		return nil, errutil.Err(err)
	}

	fg, err := util.ParseColor(jt.Graphics.Fg)
	if err != nil {
		return nil, errutil.Err(err)
	}
	bg, err := util.ParseColor(jt.Graphics.Bg)
	if err != nil {
		return nil, errutil.Err(err)
	}

	t = &Terrain{
		name: jt.Name,
	}
	t.SetPathable(jt.Pathable)
	t.SetGraphics(termbox.Cell{
		Ch: rune(jt.Graphics.Ch[0]),
		Fg: fg,
		Bg: bg,
	})
	return t, nil
}
