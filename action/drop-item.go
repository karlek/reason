package action

import (
	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/status"

	"github.com/nsf/termbox-go"
)

func DropItem(hero *beastiary.Creature, a *area.Area) bool {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()

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

func NarrativeDropItem(ch string, hero *beastiary.Creature, a *area.Area) {
	i := hero.DropItem(ch, a)
	status.Print("You dropped " + i.Name() + ".")
}

func PrintCategorizedInventory(title string, hero *beastiary.Creature) {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	if len(hero.Inventory) == 0 {
		ui.PrintInventory("You aren't carrying anything.", 0, 0, ui.WholeScreenWidth, termbox.ColorWhite+termbox.AttrBold, termbox.ColorDefault)
		termbox.PollEvent()
		return
	}
	ui.PrintInventory(title, 0, 0, ui.WholeScreenWidth, termbox.ColorWhite+termbox.AttrBold, termbox.ColorDefault)

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
			s := i.Hotkey + " - " + i.Name()
			if _, ok := categories[i.Category]; !ok {
				categories["unknown"] = append(categories["unknown"], s)
			} else {
				categories[i.Category] = append(categories[i.Category], s)
			}
		}
	}

	yOffset := 2
	xOffset := 1
	row := 0
	for cat, items := range categories {
		if len(items) == 0 {
			continue
		}
		ui.PrintInventory(properCategory[cat], 0, row+yOffset, ui.WholeScreenWidth, termbox.ColorCyan+termbox.AttrBold, termbox.ColorDefault)
		row++
		for _, itemStr := range items {
			ui.PrintInventory(itemStr, xOffset, row+yOffset, ui.WholeScreenWidth, termbox.ColorWhite, termbox.ColorDefault)
			row++
		}
		row++
	}
}
