package creature

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/mewkiz/pkg/errutil"
)

const (
	UnableToDrop = "Couldn't drop item."
)

type Inventory map[rune]item.DrawItemer

func (c *Creature) Equip(pos rune) item.Itemer {
	i := c.Inventory[pos]
	switch e := i.(type) {
	case *item.Weapon:
		c.Equipment.MainHand = e
		return i
	case *item.Ring:
		c.Equipment.Rings = append(c.Equipment.Rings, e)
		return i
	}
	return nil
}

func (c *Creature) IsEquipped(i item.Itemer) (isEquipped bool) {
	if i == nil {
		return
	}
	switch obj := i.(type) {
	case (*item.Weapon):
		isEquipped = c.Equipment.MainHand == obj || c.Equipment.OffHand == obj
	case (*item.Headgear):
		isEquipped = c.Equipment.Head == obj
	case (*item.Amulet):
		isEquipped = c.Equipment.Amulet == obj
	case (*item.Ring):
		isEquipped = inRings(obj, c.Equipment.Rings)
	case (*item.Boots):
		isEquipped = c.Equipment.Boots == obj
	case (*item.Gloves):
		isEquipped = c.Equipment.Gloves == obj
	case (*item.Chestwear):
		isEquipped = c.Equipment.Chestwear == obj
	case (*item.Legwear):
		isEquipped = c.Equipment.Legwear == obj
	}
	return
}

func inRings(needle *item.Ring, rings []*item.Ring) bool {
	for _, ring := range rings {
		if needle == ring {
			return true
		}
	}
	return false
}

func (c *Creature) Use(i item.Itemer) {
	if !item.IsUsable(i) {
		status.Print("You can't use that item!")
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

func (c *Creature) use(i item.Itemer) {
	switch i.Name() {
	case "Potion of Increase Weight":
		status.Print("You drink the potion.")
		status.Print("It tastes like metal. Your backpack feels much heavier!")
	case "Star-Eye Map":
		status.Print("You try to read the map.")
		status.Print("It's exhausting, but you know exactly how to reach your goal!")
	}
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
		c.removeRing(obj)
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
	status.Printf("You unequip %s.", i.Name())
}

func (c *Creature) removeRing(obj *item.Ring) {
	for index, ring := range c.Equipment.Rings {
		if obj == ring {
			c.Equipment.Rings[index] = nil
			return
		}
	}
}

func (c *Creature) PickUp(a *area.Area) (actionTaken bool) {
	msg := "There's no item here."
	i, err := c.pickUp(a)
	if i == nil {
		status.Print(msg)
		return false
	}

	// Print status message if hero's inventory is full.
	if c.IsHero() {
		if err != nil {
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("%c - %s picked up.", i.Hotkey(), i.String())
		}
	} else {
		msg = fmt.Sprintf("%s picked up %s.", strings.Title(c.Name()), i.String())
	}

	// If the distance to the creature is within the sight radius, print the
	// status message.
	if c.dist() <= Hero.Sight {
		status.Print(msg)
	}

	return true
}

func (c *Creature) pickUp(a *area.Area) (item.DrawItemer, error) {
	// Take the topmost item of the cell the creature is standing on.
	stk, ok := a.Items[c.Coord()]
	if !ok || stk.Len() == 0 {
		return nil, nil
	}
	i, ok := stk.Pop().(item.DrawItemer)
	if !ok {
		return nil, nil
	}

	var hotkey rune
	var err error
	// Tries to find existing stack of item and increase it's count, otherwise
	// try to add it to the inventory if it's not full.
	if hotkey, ok = c.findStack(i); ok {
		inv := c.Inventory[hotkey]
		inv.SetCount(inv.Count() + i.Count())
	} else {
		hotkey, err = c.findHotkey()
		if err != nil {
			return nil, err
		}
		c.Inventory[hotkey] = i
	}
	i.SetHotkey(hotkey)
	return i, nil
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
		status.Print(UnableToDrop)
		return
	}

	fmtStr := "%s dropped %s."
	cName := strings.Title(c.Name())
	if c.IsHero() {
		cName = "You"
	}
	iName := i.Name()
	if item.IsStackable(i) {
		iName = strconv.Itoa(i.Count()) + " " + i.Name()
	}

	if c.dist() <= Hero.Sight {
		status.Printf(fmtStr, cName, iName)
	}
}

// findStack takes an item and tries to find a stack of that item in the
// inventory. If a stack exists it returns the item slot letter and true,
// otherwise it returns empty rune and false.
func (c *Creature) findStack(i item.DrawItemer) (hotkey rune, ok bool) {
	if !item.IsStackable(i) {
		return 0x00, false
	}
	for _, v := range c.Inventory {
		if v.Name() == i.Name() {
			return v.Hotkey(), true
		}
	}
	return 0x00, false
}

// findHotkey returns the first free inventory slot or returns
// "Inventory is full." error.
func (c *Creature) findHotkey() (hotkey rune, err error) {
	for _, hotkey := range item.Positions {
		if _, ok := c.Inventory[hotkey]; !ok {
			return hotkey, nil
		}
	}
	return 0x00, errutil.Newf("Inventory is full.")
}

func (inv Inventory) FreeSlots() (num int) {
	for _, hotkey := range item.Positions {
		if _, ok := inv[hotkey]; !ok {
			num++
		}
	}
	return num
}

func (inv Inventory) UsedSlots() (num int) {
	for _, hotkey := range item.Positions {
		if _, ok := inv[hotkey]; ok {
			num++
		}
	}
	return num
}
