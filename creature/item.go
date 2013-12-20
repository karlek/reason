package creature

import (
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/mewkiz/pkg/errutil"
)

type Inventory map[string]*item.Item

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
		hotkey, err := c.findHotkey()
		if err != nil {
			status.Print(err.Error())
			return nil
		}
		i.Hotkey = hotkey
		c.Inventory[i.Hotkey] = i
	}
	return i
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

func (c *Creature) findHotkey() (hotkey string, err error) {
	for _, ch := range item.Letters {
		hotkey := string(ch)
		if _, ok := c.Inventory[hotkey]; !ok {
			return hotkey, nil
		}
	}
	return "", errutil.Newf("Inventory is full.")
}
