package creature

import (
	"fmt"

	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/nsf/termbox-go"
)

// Action performs simple AI for a creature.
func (c *Creature) Action(a *area.Area, hero *Creature) int {
	var col *area.Collision
	var err error
	if c.X() < hero.X() {
		col, err = a.MoveRight(c)
	} else if c.X() > hero.X() {
		col, err = a.MoveLeft(c)
	} else if c.Y() < hero.Y() {
		col, err = a.MoveDown(c)
	} else if c.Y() > hero.Y() {
		col, err = a.MoveUp(c)
	}
	if err != nil {
		return 0
	}
	if col == nil {
		return c.Speed
	}
	if mob, ok := col.S.(*Creature); ok {
		if mob.Name() == "hero" {
			battleNarrative(a, hero, c)
			return c.Speed
		}
		/// mobs friendly fire
		// } else {
		// battle(a, c, mob)
		// }
	}

	/// this simulates waiting.
	return c.Speed
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
		termbox.Flush()
		termbox.PollEvent()
		util.Quit()
	}
}

// Actions performs a number of actions for the creature.
func (c *Creature) Actions(turns int, a *area.Area, hero *Creature) {
	for ; turns > 0; turns-- {
		c.Action(a, hero)
	}
}
