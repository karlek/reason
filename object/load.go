package object

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/karlek/reason/util"

	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

var Objects = map[string]Object{}

// Load initializes the objects collection.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/object/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		o, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		log.Println("here!")
		Objects[o.Name()] = *o
	}
	return nil
}

// load parses a JSON data file into a go object.
func load(filename string) (o *Object, err error) {
	// jsonobject is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonobject struct {
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

	var jo jsonobject
	err = json.Unmarshal(buf, &jo)
	if err != nil {
		return nil, errutil.Err(err)
	}

	fg, err := util.ParseColor(jo.Graphics.Fg)
	if err != nil {
		return nil, errutil.Err(err)
	}
	bg, err := util.ParseColor(jo.Graphics.Bg)
	if err != nil {
		return nil, errutil.Err(err)
	}

	o = &Object{
		name: jo.Name,
	}
	o.SetPathable(jo.Pathable)
	o.SetGraphics(termbox.Cell{
		Ch: rune(jo.Graphics.Ch[0]),
		Fg: fg,
		Bg: bg,
	})
	return o, nil
}
