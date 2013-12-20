// Package creature contains information about all creatures in reason.
package creature

import (
	"github.com/karlek/worc/model"
	"github.com/nsf/termbox-go"
)

// Creature is a model with stats, name, inventory... etc.
type Creature struct {
	M         model.Model
	name      string
	MaxHp     int
	Hp        int
	Strength  int
	Sight     int
	CurSpeed  int
	Speed     int
	Inventory Inventory
}

// Name returns the name of the creature.
func (c *Creature) Name() string {
	return c.name
}

// NewX sets a new x value for the coordinate.
func (c *Creature) NewX(x int) {
	c.M.NewX(x)
}

// NewY sets a new y value for the coordinate.
func (c *Creature) NewY(y int) {
	c.M.NewY(y)
}

// IsPathable returns whether objects can be stacked ontop of this object.
func (c *Creature) IsPathable() bool {
	return c.M.IsPathable()
}

// Graphic returns the graphic data of this object.
func (c *Creature) Graphic() termbox.Cell {
	return c.M.Graphic()
}

// X returns the x value of the current coordinate.
func (c *Creature) X() int {
	return c.M.X()
}

// Y returns the y value of the current coordinate.
func (c *Creature) Y() int {
	return c.M.Y()
}
