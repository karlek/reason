package creature

import (
	"fmt"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/mewkiz/pkg/errutil"
)

type Inventory map[rune]item.Itemer

type Equipment struct {
	MainHand  *item.Weapon
	OffHand   *item.Weapon
	Head      *item.Headgear
	Amulet    *item.Amulet
	Rings     []*item.Ring
	Boots     *item.Boots
	Gloves    *item.Gloves
	Chestwear *item.Chestwear
	Legwear   *item.Legwear
}

func (c *Creature) Equip(pos rune) item.Itemer {
	i := c.Inventory[pos]
	switch e := i.(type) {
	case *item.Weapon:
		c.Equipment.MainHand = e
		return i
	}
	status.Print(fmt.Sprintf("i: %T", i))
	return nil
}

func (c *Creature) PickUp(a *area.Area) item.Itemer {
	s, ok := a.Items[c.Coord()]
	if !ok {
		return nil
	}
	i := s.Pop()
	if i == nil {
		return nil
	}

	if hotkey, ok := c.findStack(i); ok {
		i.SetHotkey(hotkey)
		c.Inventory[i.Hotkey()].IncCount(i.Count())
	} else {
		hotkey, err := c.findHotkey()
		if err != nil {
			status.Print(err.Error())
			return nil
		}
		i.SetHotkey(hotkey)
		c.Inventory[hotkey] = i
	}
	return i
}

func (c *Creature) DropItem(pos rune, a *area.Area) item.Itemer {
	i := c.Inventory[pos]
	delete(c.Inventory, pos)

	cor := c.Coord()
	if a.Items[cor] == nil {
		a.Items[cor] = new(area.Stack)
	}
	a.Items[cor].Push(i)

	return i
}

// findStack takes an item and tries to find a stack of that item in the
// inventory. If a stack exists it returns the item slot letter and true,
// otherwise it returns empty rune and false.
func (c *Creature) findStack(i item.Itemer) (hotkey rune, ok bool) {
	if !i.IsStackable() {
		return "", false
	}
	for _, v := range c.Inventory {
		if v.Name() == i.Name() {
			return v.Hotkey(), true
		}
	}
	return "", false
}

// findHotkey returns the first free inventory slot or returns
// "Inventory is full." error.
func (c *Creature) findHotkey() (hotkey rune, err error) {
	for _, hotkey := range item.Positions {
		if _, ok := c.Inventory[hotkey]; !ok {
			return hotkey, nil
		}
	}
	return "", errutil.Newf("Inventory is full.")
}

func (inv Inventory) getFreeSlots() (num int) {
	for _, hotkey := range item.Positions {
		if _, ok := inv[hotkey]; !ok {
			num++
		}
	}
	return num
}

func (inv Inventory) GetUsedSlots() (num int) {
	for _, hotkey := range item.Positions {
		if _, ok := inv[hotkey]; ok {
			num++
		}
	}
	return num
}
