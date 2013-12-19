// Package beastiary contains information about all creatures in reason.
package beastiary

import (
	"fmt"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/model"
	"github.com/nsf/termbox-go"
)

type Inventory map[string]*item.Item

// Creature is a model with stats, name, inventory... etc.
type Creature struct {
	M         model.Model
	N         string
	MaxHp     int
	Hp        int
	Strength  int
	Sight     int
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

	if hotkey, ok := c.findStack(i); ok {
		i.Hotkey = hotkey
		c.Inventory[i.Hotkey].Num += i.Num
	} else {
		i.Hotkey = c.findHotkey()
		c.Inventory[i.Hotkey] = i
	}
	return i
}

func (c *Creature) findStack(i *item.Item) (hotkey string, ok bool) {
	if !i.IsStackable() {
		return "", false
	}
	for _, v := range c.Inventory {
		if v.Name() == i.Name() {
			return v.Hotkey, true
		}
	}
	return "", false
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
		/// mobs friendly fire
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
}

func battleNarrative(a *area.Area, hero *Creature, attacker *Creature) {
	hero.Hp -= attacker.Strength
	status.Print(fmt.Sprintf("You take %d damage from %s!", attacker.Strength, attacker.Name()))
	if hero.Hp <= 0 {
		hero.DrawFOV(a)
		status.Print("You die. Press any key to quit.")
		termbox.PollEvent()
		util.Quit()
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
