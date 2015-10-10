// Package creature contains information about all creatures in reason.
package creature

import (
	"math"

	"github.com/karlek/reason/creature/equipment"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/model"
)

// Creature is a model with stats, name, inventory... etc.
type Creature struct {
	name.Namer
	model.Model
	name         string
	MaxHp        int
	Hp           int
	Strength     int
	Sight        int
	Speed        int
	Regeneration int
	RegCounter   int
	Level        int
	Experience   int
	Inventory    Inventory
	Equipment    equipment.Equipment
}

// Name returns the name of the creature.
func (c *Creature) Name() string {
	return c.name
}

func (c *Creature) Corpse() item.DrawItemer {
	return item.Items[c.Name()+" corpse"]
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

func (c *Creature) MonstersInRange(a *area.Area) []*Creature {
	var monsters []*Creature

	cs := c.FOV(a)
	for p := range cs {
		if monst, ok := a.Monsters[p]; ok {
			// Ignore hero.
			if monst == nil || monst.Name() == "hero" {
				continue
			}
			monsters = append(monsters, monst.(*Creature))
		}
	}
	return monsters
}

func (c *Creature) Reg() {
	if c.Hp == c.MaxHp {
		return
	}
	// Unit doesn't have regeneration.
	if c.Regeneration == 0 {
		return
	}
	c.RegCounter++
	if c.RegCounter == c.Regeneration {
		c.Hp++
		c.RegCounter = 0
	}
}
