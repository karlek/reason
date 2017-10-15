package creature

import (
	"github.com/karlek/reason/item"

	"github.com/karlek/reason/area"
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

func (c *Creature) removeRing(obj *item.Ring) {
	for index, ring := range c.Equipment.Rings {
		if obj == ring {
			c.Equipment.Rings[index] = nil
			return
		}
	}
}

func (c *Creature) PickUp(a *area.Area) (item.DrawItemer, error) {
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
