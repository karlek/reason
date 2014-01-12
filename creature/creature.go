// Package creature contains information about all creatures in reason.
package creature

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/model"
)

// Creature is a model with stats, name, inventory... etc.
type Creature struct {
	name.Namer
	model.Model
	name      string
	MaxHp     int
	Hp        int
	Strength  int
	Sight     int
	Speed     int
	Inventory Inventory
	Equipment Equipment
}

// Name returns the name of the creature.
func (c *Creature) Name() string {
	return c.name
}
