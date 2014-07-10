// Package creature contains information about all creatures in reason.
package creature

import (
	"math"

	"github.com/karlek/reason/creature/equipment"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
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

	// Inclusive hero's square so it's from hero's eyes.
	radius := c.Sight
	for x := c.X() - radius; x <= c.X()+radius; x++ {
		for y := c.Y() - radius; y <= c.Y()+radius; y++ {
			// Discriminate coordinates which are out of bounds.
			if !a.ExistsXY(x, y) {
				continue
			}

			// Distance between creature x and y coordinates and sight radius.
			dx := float64(x - c.X())
			dy := float64(y - c.Y())

			// Distance between creature and sight radius.
			dist := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

			// Discriminate coordinates which are outside of the circle.
			if dist > float64(radius) {
				continue
			}

			//
			cor := coord.Coord{x, y}
			if monst, ok := a.Monsters[cor]; ok {
				if monst == nil || monst.Name() == "hero" {
					continue
				}
				monsters = append(monsters, monst.(*Creature))
			}
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
