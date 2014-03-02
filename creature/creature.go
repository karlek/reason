// Package creature contains information about all creatures in reason.
package creature

import (
	"math"

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

// IsHero compares creature c with Hero.
func (c Creature) IsHero() bool {
	return c.Name() == Hero.Name()
}

// dist returns the distance between creature c and the hero.
func (c Creature) dist() int {
	x := math.Pow(float64(Hero.X()-c.X()), 2)
	y := math.Pow(float64(Hero.Y()-c.Y()), 2)

	return int(math.Sqrt(x + y))
}
