package action

import (
	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"

	"github.com/nsf/termbox-go"
)

func ShowInventory(a *area.Area, hero *beastiary.Creature) bool {
	PrintCategorizedInventory("Inventory: currentWeight/maxPossibleWeight (usedSlots/totalSlots)", hero)
	if len(hero.Inventory) == 0 {
		return false
	}

inventoryLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == ui.CancelKey {
				break inventoryLoop
			}

			hotkey := string(ev.Ch)
			if i, ok := hero.Inventory[hotkey]; ok {
				if actionTaken := ShowItemDetails(i, hero, a); actionTaken {
					return true
				} else {
					PrintCategorizedInventory("Inventory: currentWeight/maxPossibleWeight (usedSlots/totalSlots)", hero)
				}
			}
		}
	}
	return false
}

func ShowItemDetails(i *item.Item, hero *beastiary.Creature, a *area.Area) bool {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	termbox.Flush()
	rows := 0
	ui.PrintInventory(i.Hotkey+" - "+i.Name(), 0, 0, ui.WholeScreenWidth, termbox.ColorWhite, termbox.ColorDefault)
	rows += 2
	ui.PrintInventory(i.Description, 0, rows, ui.WholeScreenWidth, termbox.ColorWhite, termbox.ColorDefault)
	rows += 2
	ui.PrintInventory("You can (d)rop this item.", 0, rows, ui.WholeScreenWidth, termbox.ColorCyan, termbox.ColorDefault)

itemDetailLoop:
	for {
		switch detailsEvent := termbox.PollEvent(); detailsEvent.Type {
		case termbox.EventKey:
			if detailsEvent.Key == ui.CancelKey {
				break itemDetailLoop
			}

			itemAction := string(detailsEvent.Ch)
			if itemAction == string(ui.DropItemKey) {
				NarrativeDropItem(i.Hotkey, hero, a)
			}
			return true
		}
	}
	return false
}
