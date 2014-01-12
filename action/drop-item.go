package action

import (
	"strconv"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"

	"github.com/nsf/termbox-go"
)

func DropItem(a *area.Area, hero *creature.Creature) bool {
	PrintCategorizedInventory("Drop Item: currentWeight/maxPossibleWeight (usedSlots/totalSlots)", hero)

	if len(hero.Inventory) == 0 {
		return false
	}

dropItemLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == ui.CancelKey {
				break dropItemLoop
			}
			hotkey := string(ev.Ch)
			if _, ok := hero.Inventory[hotkey]; ok {
				NarrativeDropItem(hotkey, hero, a)
				return true
			}
		}
	}
	return false
}

func NarrativeDropItem(ch string, hero *creature.Creature, a *area.Area) {
	i := hero.DropItem(ch, a)
	if i == nil {
		status.Print("I failed :(")
	}
	s := "You dropped "
	if i.IsStackable() {
		s += strconv.Itoa(i.GetNum()) + " " + i.Name()
	} else {
		s += i.Name()
	}
	s += "."
	status.Print(s)
}

func NarrativeEquip(ch string, hero *creature.Creature) {
	i := hero.Equip(ch)
	if i == nil {
		status.Print("That item can't be equipped.")
		return
	}

	s := "You equipped "
	if i.IsStackable() {
		s += strconv.Itoa(i.GetNum()) + " " + i.Name()
	} else {
		s += i.Name()
	}
	s += "."
	status.Print(s)
}

func PrintCategorizedInventory(title string, hero *creature.Creature) {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	if len(hero.Inventory) == 0 {
		ui.PrintInventory("You aren't carrying anything.", 0, 0, ui.Whole.Width, termbox.ColorWhite+termbox.AttrBold, termbox.ColorDefault)
		termbox.Flush()
		termbox.PollEvent()
		return
	}
	ui.PrintInventory(title, 0, 0, ui.Whole.Width, termbox.ColorWhite+termbox.AttrBold, termbox.ColorDefault)

	// item category is the key to which plural category it will be placed under.
	var categories = map[string][]string{
		"tool":   {},
		"armour": {},
		"ring":   {},
		"weapon": {},
	}

	var properCategory = map[string]string{
		"tool":    "Tools",
		"armour":  "Armours",
		"ring":    "Rings",
		"weapon":  "Weapons",
		"unknown": "Unknown",
	}

	// We do this the quirky we to get the inventory list sorted.
	for _, ch := range item.Letters {
		if i, ok := hero.Inventory[string(ch)]; ok {
			var s string
			if i.IsStackable() {
				s = i.GetHotkey() + " - " + strconv.Itoa(i.GetNum()) + " " + i.Name()
			} else {
				s = i.GetHotkey() + " - " + i.Name()
			}
			if i.IsEquipable() {
				if hero.Equipment.MainHand != nil {
					if hero.Equipment.MainHand.GetHotkey() == i.GetHotkey() {
						s += " (equipped)"
					}
				}
			}
			if _, ok := categories[i.GetCategory()]; !ok {
				categories["unknown"] = append(categories["unknown"], s)
			} else {
				categories[i.GetCategory()] = append(categories[i.GetCategory()], s)
			}
		}
	}

	yOffset := 2
	xOffset := 1
	row := 0

	/// Make inventory screen in ui
	for cat, items := range categories {
		if len(items) == 0 {
			continue
		}
		ui.PrintInventory(properCategory[cat], 0, row+yOffset, ui.Whole.Width, termbox.ColorCyan+termbox.AttrBold, termbox.ColorDefault)
		row++
		for _, itemStr := range items {
			ui.PrintInventory(itemStr, xOffset, row+yOffset, ui.Whole.Width, termbox.ColorWhite, termbox.ColorDefault)
			row++
		}
		row++
	}
	termbox.Flush()
}
