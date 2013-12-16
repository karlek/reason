// Package beastiary contains information about all creatures in reason.
package beastiary

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/object"
	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

// Creatures is a map where names of the creature is the key mapping to that
// creature object.
var Creatures = map[string]Creature{}

// Creature is an object with a name.
type Creature struct {
	O        object.Object
	N        string
	MaxHp    int
	Hp       int
	Strength int
	CurSpeed float64
	Speed    float64
}

// action performs simple AI for a creature.
func (c *Creature) action(a *area.Area, hero *Creature) {
	i := randInt(0, 1000)
	var col area.Stackable
	switch {
	case i == 999:
		status.Print("Fresh grass, yum!")
	case i%4 == 0:
		col = a.MoveUp(c)
	case i%4 == 1:
		col = a.MoveDown(c)
	case i%4 == 2:
		col = a.MoveLeft(c)
	case i%4 == 3:
		col = a.MoveRight(c)
	default:
		return
	}

	if mob, ok := col.(*Creature); ok {
		if mob.Name() == "hero" {
			battleNarrative(a, hero, c)
		} else {
			battle(a, c, mob)
		}
	}
}

// battle between to non player characters.
func battle(a *area.Area, defender *Creature, attacker *Creature) {
	defender.Hp -= attacker.Strength
	if defender.Hp <= 0 {
		a.Objects[coord.Coord{defender.X(), defender.Y()}].Pop()
	}
	a.Draw()
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

// randInt is used by the debug function GenArea.
func randInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
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
