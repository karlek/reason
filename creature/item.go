package creature

import (
	"fmt"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/mewkiz/pkg/errutil"
)

type Inventory map[string]*item.Item

type Equipment struct {
	// MainHand *item.Weapon
	// // OffHand   item.Weapon
	// Head      *item.Headgear
	// Amulet    *item.Amulet
	// Rings     []*item.Ring
	// Boots     *item.Boots
	// Gloves    *item.Gloves
	// Chestwear *item.Chestwear
	// Legwear   *item.Legwear
	MainHand *item.Item
	// OffHand   item.Item
	Head      *item.Item
	Amulet    *item.Item
	Rings     []*item.Item
	Boots     *item.Item
	Gloves    *item.Item
	Chestwear *item.Item
	Legwear   *item.Item
}

func (c *Creature) Equip(ch string) *item.Item {
	i := c.Inventory[ch]
	// switch e := i.(type) {
	// case (*item.Weapon):
	// c.Equipment.MainHand = e
	// return i
	// case (*item.Item):
	if i.GetCategory() != "weapon" {
		status.Print(fmt.Sprintf("not a weapon: %T", i))
		return nil
	}
	c.Equipment.MainHand = i
	return i
	// default:
	// status.Print(fmt.Sprintf("e: %T", e))
	// return nil
	// }
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
	i, ok := w.(item.Item)
	if !ok {
		return nil
	}

	if hotkey, ok := c.findStack(i); ok {
		i.Hotkey = hotkey
		c.Inventory[i.Hotkey].IncreaseNum(i.Num)
	} else {
		hotkey, err := c.findHotkey()
		if err != nil {
			status.Print(err.Error())
			return nil
		}
		i.Hotkey = hotkey
		c.Inventory[i.Hotkey] = &i
	}
	return &i
}

func (c *Creature) DropItem(ch string, a *area.Area) *item.Item {
	i := c.Inventory[ch]
	delete(c.Inventory, ch)

	cor := coord.Coord{c.X(), c.Y()}
	if a.Items[cor] == nil {
		a.Items[cor] = new(area.Stack)
	}
	a.Items[cor].Push(*i)

	return i
}

// findStack takes an item and tries to find a stack of that item in the
// inventory. If a stack exists it returns the item slot letter and true,
// otherwise it returns empty string and false.
func (c *Creature) findStack(i item.Item) (hotkey string, ok bool) {
	if !i.IsStackable() {
		return "", false
	}
	for _, v := range c.Inventory {
		if v.Name() == i.Name() {
			return v.GetHotkey(), true
		}
	}
	return "", false
}

// findHotkey returns the first free inventory slot or returns
// "Inventory is full." error.
func (c *Creature) findHotkey() (hotkey string, err error) {
	for _, ch := range item.Letters {
		hotkey := string(ch)
		if _, ok := c.Inventory[hotkey]; !ok {
			return hotkey, nil
		}
	}
	return "", errutil.Newf("Inventory is full.")
}

func (inv Inventory) getFreeSlots() (num int) {
	for _, ch := range item.Letters {
		hotkey := string(ch)
		if _, ok := inv[hotkey]; !ok {
			num++
		}
	}
	return num
}

func (inv Inventory) GetUsedSlots() (num int) {
	for _, ch := range item.Letters {
		hotkey := string(ch)
		if _, ok := inv[hotkey]; ok {
			num++
		}
	}
	return num
}
