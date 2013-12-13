// Package fauna adds objects which are stationary doodads.
package fauna

import (
	"github.com/karlek/worc/area"
	"github.com/karlek/worc/object"
	"github.com/nsf/termbox-go"
)

// Doodad is non-moveable, but drawable object.
type Doodad struct {
	area.Stackable
	O object.Object
	N string
}

// IsStackable returns whether objects can be stacked ontop of this object.
func (d Doodad) IsStackable() bool {
	return d.O.IsStackable()
}

// Graphic is needed for draw.Drawable.
func (d Doodad) Graphic() termbox.Cell {
	return d.O.Graphic()
}

// Name returns the name of the creature.
func (d Doodad) Name() string {
	return d.N
}

// Wall is a white '#', which can not be walked over.
var Wall = Doodad{
	O: object.Object{
		G: termbox.Cell{
			Ch: '#',
			Fg: termbox.ColorWhite + termbox.AttrBold,
		},
		Stackable: false,
	},
	N: "a wall",
}

// Soil is a "yellow" '.', which can be walked over.
var Soil = Doodad{
	O: object.Object{
		G: termbox.Cell{
			Ch: '.',
			Fg: termbox.ColorYellow,
		},
		Stackable: true,
	},
	N: "soil",
}

// Water is a blue '~', which can be walked over.
var Water = Doodad{
	O: object.Object{
		G: termbox.Cell{
			Ch: '~',
			Fg: termbox.ColorBlue,
		},
		Stackable: true,
	},
	N: "water",
}
