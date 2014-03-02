// Package terrain adds functionality for stationary terrain.
package terrain

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/model"
)

// Terrain is non-moveable, but drawable object.
type Terrain struct {
	model.Model
	name.Namer
	name       string
	IsExplored bool
}

// Name returns the name of the creature.
func (t Terrain) Name() string {
	return t.name
}

func New(t Terrain) *Terrain {
	return &t
}
