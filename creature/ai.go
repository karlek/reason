package creature

import (
	"fmt"
	"strings"

	// "github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/nsf/termbox-go"
)

// Action performs simple AI for a creature.
func (c *Creature) Action(a *area.Area) int {
	if i := a.Items[c.Coord()].Peek(); i != nil {
		c.PickUp(a)
		return c.Speed
	}
	//  else if len(c.Inventory) > 1 {
	// 	for _, pos := range item.Positions {
	// 		if _, ok := c.Inventory[pos]; ok {
	// 			c.DropItem(pos, a)
	// 			return c.Speed
	// 		}
	// 	}
	// }
	var col *area.Collision
	var err error
	if c.X() < Hero.X() {
		col, err = a.MoveRight(c)
	} else if c.X() > Hero.X() {
		col, err = a.MoveLeft(c)
	} else if c.Y() < Hero.Y() {
		col, err = a.MoveDown(c)
	} else if c.Y() > Hero.Y() {
		col, err = a.MoveUp(c)
	}
	if err != nil {
		return 0
	}
	if col == nil {
		return c.Speed
	}
	if mob, ok := col.S.(*Creature); ok {
		if mob.IsHero() {
			c.Battle(mob, a)
			return c.Speed
		}
	}

	// If all fails, creature waits.
	return c.Speed
}

func (attacker *Creature) Battle(defender *Creature, a *area.Area) {
	var s string
	if defender.IsHero() {
		s = fmt.Sprintf("You take %d damage from %s!", attacker.Strength, attacker.Name())
	} else if attacker.IsHero() {
		s = fmt.Sprintf("You inflict %d damage to %s!", attacker.Strength, defender.Name())
	} else {
		s = fmt.Sprintf("%s takes %d damage from %s!", strings.Title(defender.Name()), attacker.Strength, attacker.Name())
	}

	defender.Hp -= attacker.Strength
	if defender.Hp <= 0 {
		if defender.IsHero() {
			Hero.DrawFOV(a)
			status.Print(s)
			status.Print("You die. Press any key to quit.")
			termbox.Flush()
			termbox.PollEvent()
			util.Quit()
		} else if attacker.IsHero() {
			s += fmt.Sprintf(" You killed %s!", defender.Name())
		}
		a.Monsters[coord.Coord{defender.X(), defender.Y()}] = nil
		_, ok := a.Items[defender.Coord()]
		if !ok {
			a.Items[defender.Coord()] = new(area.Stack)
		}
		for _, i := range defender.Inventory {
			a.Items[defender.Coord()].Push(i)
		}
	}
	if attacker.dist() <= Hero.Sight {
		status.Print(s)
	}
}
