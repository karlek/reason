// Package beastiary contains information about all creatures in reason.
package beastiary

import (
	"fmt"
	"os"
	"time"

	"github.com/karlek/reason/item"
	// "github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/object"
	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

type Inventory map[string]*item.Item

// Creature is an object with a name.
type Creature struct {
	O         object.Object
	N         string
	MaxHp     int
	Hp        int
	Strength  int
	CurSpeed  float64
	Speed     float64
	Inventory Inventory
}

func (c *Creature) PickUp(a *area.Area) *item.Item {
	x, y := c.X(), c.Y()
	cor := coord.Coord{x, y}

	s, ok := a.Items[cor]
	if !ok {
		return nil
	}
	w := s.Pop()
	if w == nil {
		return nil
	}
	i, ok := w.(*item.Item)
	if !ok {
		return nil
	}

	i.Hotkey = c.findHotkey()
	c.Inventory[i.Hotkey] = i
	return i
}

func (c *Creature) DropItem(ch string, a *area.Area) *item.Item {
	i := c.Inventory[ch]
	delete(c.Inventory, ch)

	cor := coord.Coord{c.X(), c.Y()}
	if a.Items[cor] == nil {
		a.Items[cor] = new(area.Stack)
	}
	a.Items[cor].Push(i)

	return i
}

/// can overflow!!
func (c *Creature) findHotkey() string {
	for _, ch := range item.Letters {
		hotkey := string(ch)
		if _, ok := c.Inventory[hotkey]; !ok {
			return hotkey
		}
	}
	return ""
}

// action performs simple AI for a creature.
func (c *Creature) action(a *area.Area, hero *Creature) {
	var col *area.Collision
	if c.X() < hero.X() {
		col = a.MoveRight(c)
	} else if c.X() > hero.X() {
		col = a.MoveLeft(c)
	} else if c.Y() < hero.Y() {
		col = a.MoveDown(c)
	} else if c.Y() > hero.Y() {
		col = a.MoveUp(c)
	}

	// i := util.RandInt(0, 1000)
	// switch {
	// case i == 999:
	// 	status.Print("Fresh grass, yum!")
	// case i%4 == 0:
	// 	col = a.MoveUp(c)
	// case i%4 == 1:
	// 	col = a.MoveDown(c)
	// case i%4 == 2:
	// 	col = a.MoveLeft(c)
	// case i%4 == 3:
	// 	col = a.MoveRight(c)
	// default:
	// 	return
	// }
	if col == nil {
		return
	}
	if mob, ok := col.S.(*Creature); ok {
		if mob.Name() == "hero" {
			battleNarrative(a, hero, c)
		}
		// } else {
		// battle(a, c, mob)
		// }
	}
}

// battle between to non player characters.
func battle(a *area.Area, defender *Creature, attacker *Creature) {
	defender.Hp -= attacker.Strength
	if defender.Hp <= 0 {
		a.Monsters[coord.Coord{defender.X(), defender.Y()}] = nil
	}
	a.ReDraw(defender.X(), defender.Y())
}

///
func battleNarrative(a *area.Area, hero *Creature, attacker *Creature) {
	hero.Hp -= attacker.Strength
	status.Print(fmt.Sprintf("You take %d damage from %s!", attacker.Strength, attacker.Name()))
	if hero.Hp <= 0 {
		status.Print("You die. Press any key to quit.")
		termbox.PollEvent()
		time.Sleep(100)
		termbox.Close()
		os.Exit(0)
	}
}

// Actions performs a number of actions for the creature.
func (c *Creature) Actions(turns int, a *area.Area, hero *Creature) {
	for ; turns > 0; turns-- {
		c.action(a, hero)
	}
}

// Name returns the name of the creature.
func (c *Creature) Name() string {
	return c.N
}

// NewX sets a new x value for the coordinate.
func (c *Creature) NewX(x int) {
	c.O.NewX(x)
}

// NewY sets a new y value for the coordinate.
func (c *Creature) NewY(y int) {
	c.O.NewY(y)
}

// IsStackable returns whether objects can be stacked ontop of this object.
func (c *Creature) IsStackable() bool {
	return c.O.IsStackable()
}

// Graphic returns the graphic data of this object.
func (c *Creature) Graphic() termbox.Cell {
	return c.O.Graphic()
}

// X returns the x value of the current coordinate.
func (c *Creature) X() int {
	return c.O.X()
}

// Y returns the y value of the current coordinate.
func (c *Creature) Y() int {
	return c.O.Y()
}
