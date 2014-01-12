// Package fauna adds objects which are stationary doodads.
package fauna

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/model"
)

// Doodad is non-moveable, but drawable object.
type Doodad struct {
	model.Model
	name.Namer
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

// Name returns the name of the creature.
func (d Doodad) Name() string {
	return d.name
}
