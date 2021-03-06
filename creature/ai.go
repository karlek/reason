package creature

import (
	"fmt"
	// "math/rand"
	"strings"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/ui/text"
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
	if c.Equipment.MainHand == nil && len(c.Inventory) > 0 {
		for _, pos := range item.Positions {
			if i, ok := c.Inventory[pos]; ok {
				if !item.IsEquipable(i) {
					break
				}
				if i.Name() != "Iron Sword" {
					break
				}
				c.Equip(pos)
				return c.Speed
			}
		}
	}

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
		// log.Println("err / collide?")
		return c.Speed
		// return 0
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

func (c *Creature) power() int {
	return c.Strength + c.Equipment.Power()
}

func (c *Creature) defense() int {
	return c.Equipment.Defense()
}

func (attacker *Creature) isHit(defender *Creature) bool {
	return true
}

func (attacker *Creature) Battle(defender *Creature, a *area.Area) {
	t := text.New("", termbox.ColorWhite)
	if attacker.isHit(defender) {
		t.Text = attacker.damage(defender, a)
	} else {
		t.Text = attacker.hitFail(defender)
		t.Attr = termbox.ColorBlack
	}
	if attacker.dist() <= Hero.Sight {
		status.PrintTextln(t)
	}
}

func (attacker *Creature) hitFail(defender *Creature) (s string) {
	if defender.IsHero() {
		s = fmt.Sprintf("%s misses you!", strings.Title(attacker.Name()))
	} else if attacker.IsHero() {
		s = fmt.Sprint("You miss!")
	} else {
		s = fmt.Sprintf("%s misses!", strings.Title(defender.Name()))
	}
	return s
}

func (attacker *Creature) damage(defender *Creature, a *area.Area) (s string) {
	lossOfHp := attacker.power() - defender.defense()
	if lossOfHp < 0 {
		lossOfHp = 0
	}
	if defender.IsHero() {
		s = fmt.Sprintf("You take %d damage from %s!", lossOfHp, attacker.Name())
	} else if attacker.IsHero() {
		s = fmt.Sprintf("You inflict %d damage to %s!", lossOfHp, defender.Name())
	} else {
		s = fmt.Sprintf("%s takes %d damage from %s!", strings.Title(defender.Name()), lossOfHp, attacker.Name())
	}

	defender.Hp -= lossOfHp
	if defender.Hp <= 0 {
		if defender.IsHero() {
			Hero.DrawFOV(a)
			status.Println(s, termbox.ColorWhite)
			status.Println("You die. Press any key to quit.", termbox.ColorWhite)
			status.Update()
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
		a.Items[defender.Coord()].Push(defender.Corpse())
	}
	return s
}

// func (attacker *Creature) battle(defender *Creature, a *area.Area) {

// }
