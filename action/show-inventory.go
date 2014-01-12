package action

import (
	"fmt"
	"strconv"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"

	"github.com/nsf/termbox-go"
)

func ShowInventory(a *area.Area, hero *creature.Creature) bool {
	weightStr := "Current Weight/Max Possible Weight To Carry kg"
	slotsStr := strconv.Itoa(hero.Inventory.GetUsedSlots()) + "/" + strconv.Itoa(len(item.Letters))
	inventoryDetailStr := fmt.Sprintf("Inventory: %s (%s)", weightStr, slotsStr)
	PrintCategorizedInventory(inventoryDetailStr, hero)
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
					PrintCategorizedInventory(inventoryDetailStr, hero)
				}
			}
		}
	}
	return false
}

func ShowItemDetails(i item.Itemer, hero *creature.Creature, a *area.Area) bool {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	rows := 0
	msg := MakePrintableString(i.GetHotkey() + " - " + i.Name())
	PrintLong(msg, rows)
	rows += len(msg) + 1
	msg = MakePrintableString(i.GetDescription())
	PrintLong(msg, rows)
	rows += len(msg) + 1

	str := "You can (d)rop this item."
	if i.IsEquipable() {
		str += " You can (e)quip this " + i.GetCategory() + "."
	}
	msg = MakePrintableString(str)
	PrintLongCyan(msg, rows)

	termbox.Flush()

itemDetailLoop:
	for {
		switch detailsEvent := termbox.PollEvent(); detailsEvent.Type {
		case termbox.EventKey:
			if detailsEvent.Key == ui.CancelKey {
				break itemDetailLoop
			}

			itemAction := string(detailsEvent.Ch)
			if itemAction == string(ui.DropItemKey) {
				NarrativeDropItem(i.GetHotkey(), hero, a)
			}
			if itemAction == string(ui.EquipItemKey) {
				NarrativeEquip(i.GetHotkey(), hero)
			}
			return true
		}
	}
	return false
}

// Print writes a string to the status buffer.
func MakePrintableString(str string) []string {
	var msg []string
	for {
		if len(str) < ui.Inventory.Width {
			msg = append(msg, str)
			break
		}
		strLen := ui.Inventory.Width

		msg = append(msg, str[:strLen])
		str = str[strLen:]
	}
	return msg
}

func PrintLong(msg []string, yoffset int) {
	for y, m := range msg {
		ui.PrintInventory(m, ui.Inventory.XOffset, ui.Inventory.YOffset+y+yoffset, ui.Inventory.Width, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
}

func PrintLongCyan(msg []string, yoffset int) {
	for y, m := range msg {
		ui.PrintInventory(m, ui.Inventory.XOffset, ui.Inventory.YOffset+y+yoffset, ui.Inventory.Width, termbox.ColorCyan, termbox.ColorDefault)
	}
	termbox.Flush()
}
