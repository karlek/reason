// Package fauna adds objects which are stationary doodads.
package fauna

import (
	"github.com/karlek/worc/area"
	"github.com/karlek/worc/model"
	"github.com/nsf/termbox-go"
)

// Doodad is non-moveable, but drawable object.
type Doodad struct {
	area.Tile
	M        model.Model
	name     string
	Explored bool
	BlockLOS bool
}

func (d Doodad) IsExplored() bool {
	return d.Explored
}

func (d Doodad) IsBlockingLineOfSight() bool {
	return d.BlockLOS
}

/// needs major(!) rewrite if SetExplored takes pointer reciever.
func (d Doodad) SetExplored(mode bool) {
	d.Explored = mode
}

// IsPathable returns whether objects can be stacked ontop of this object.
func (d Doodad) IsPathable() bool {
	return d.M.IsPathable()
}

// Graphic is needed for draw.Drawable.
func (d Doodad) Graphic() termbox.Cell {
	return d.M.Graphic()
}

// Name returns the name of the creature.
func (d Doodad) Name() string {
	return d.name
}
