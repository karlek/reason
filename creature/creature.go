// Package creature contains information about all creatures in reason.
package creature

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/acsellers/inflections"
	"github.com/karlek/reason/area"
	"github.com/karlek/reason/creature/equipment"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/model"
	"github.com/karlek/reason/name"
	"github.com/karlek/reason/ui/status"
	"github.com/nsf/termbox-go"
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

func (c *Creature) Use(i item.Itemer) {
	if !item.IsUsable(i) {
		status.Println("You can't use that item!", termbox.ColorRed+termbox.AttrBold)
		return
	}
	if !item.IsPermanent(i) {
		if i.Count() > 1 {
			i.SetCount(i.Count() - 1)
		} else {
			delete(c.Inventory, i.Hotkey())
		}
	}
	c.use(i)
}

func (c *Creature) UnEquip(i item.Itemer) {
	if !c.IsEquipped(i) {
		return
	}

	switch obj := i.(type) {
	case (*item.Weapon):
		if c.Equipment.MainHand == obj {
			c.Equipment.MainHand = nil
		}
		if c.Equipment.OffHand == obj {
			c.Equipment.OffHand = nil
		}
	case (*item.Headgear):
		if c.Equipment.Head == obj {
			c.Equipment.Head = nil
		}
	case (*item.Amulet):
		if c.Equipment.Amulet == obj {
			c.Equipment.Amulet = nil
		}
	case (*item.Ring):
		status.Println("Not implemented", termbox.AttrBold+termbox.ColorRed)
		return
	case (*item.Boots):
		if c.Equipment.Boots == obj {
			c.Equipment.Boots = nil
		}
	case (*item.Gloves):
		if c.Equipment.Gloves == obj {
			c.Equipment.Gloves = nil
		}
	case (*item.Chestwear):
		if c.Equipment.Chestwear == obj {
			c.Equipment.Chestwear = nil
		}
	case (*item.Legwear):
		if c.Equipment.Legwear == obj {
			c.Equipment.Legwear = nil
		}
	}
	status.Println(fmt.Sprintf("You unequip %s.", i.Name()), termbox.ColorWhite)
}

func (c *Creature) DropItem(pos rune, a *area.Area) {
	i := c.Inventory[pos]
	c.UnEquip(i)
	delete(c.Inventory, pos)

	cor := c.Coord()
	if a.Items[cor] == nil {
		a.Items[cor] = new(area.Stack)
	}
	a.Items[cor].Push(i)

	// If the item couldn't be dropped (cursed for example), print unable to
	// drop message.
	if i == nil {
		status.Println(UnableToDrop, termbox.ColorRed+termbox.AttrBold)
		return
	}

	fmtStr := "%s dropped %s."
	cName := strings.Title(c.Name())
	if c.IsHero() {
		cName = "You"
	}
	iName := i.Name()
	if item.IsStackable(i) {
		name := i.Name()
		if i.Count() > 1 {
			name = inflections.Pluralize(name)
		}
		iName = strconv.Itoa(i.Count()) + " " + name
	}

	if c.Dist() <= Hero.Sight {
		status.Println(fmt.Sprintf(fmtStr, cName, iName), termbox.ColorWhite)
	}
}

func (c *Creature) use(i item.Itemer) {
	switch i.(type) {
	case *item.Potion:
		status.Println("You drink the potion.", termbox.ColorWhite)
	case *item.Tool:
		switch i.Name() {
		case "Star-Eye Map":
			status.Println("You try to read the map.", termbox.ColorWhite)
		}
	}
	status.Print(i.UseText(), termbox.ColorWhite)
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
func (c Creature) Dist() int {
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
