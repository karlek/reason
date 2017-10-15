package inventory

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/item/effect"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/text"

	"github.com/karlek/reason/area"

	"github.com/mewkiz/pkg/stringsutil"
	"github.com/nsf/termbox-go"
)

const (
	InvTitleFmt = "Inventory: %s (%s)"
	WeightStr   = "Current Weight/Max Possible Weight To Carry kg"
	dropAction  = "You can (d)rop this item."
)

func SlotsString() string {
	return strconv.Itoa(creature.Hero.Inventory.UsedSlots()) + "/" + strconv.Itoa(len(item.Positions))
}

func Show(a *area.Area) bool {
	CategorizedInv(fmt.Sprintf(InvTitleFmt, WeightStr, SlotsString()))
	if len(creature.Hero.Inventory) == 0 {
		return false
	}

inventoryLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == ui.CancelKey {
				break inventoryLoop
			}

			if i, ok := creature.Hero.Inventory[ev.Ch]; ok {
				if actionTaken := ShowItemDetails(i, a); actionTaken {
					return true
				} else {
					CategorizedInv(fmt.Sprintf(InvTitleFmt, WeightStr, SlotsString()))
				}
			}
		}
	}
	return false
}

func StringEffects(effs map[effect.Type]effect.Magnitude) string {
	var effStr string
	for t, mag := range effs {
		effStr += fmt.Sprintf("+%d to %s\n", mag, t)
	}
	return effStr
}

func ShowItemDetails(i item.Itemer, a *area.Area) bool {
	ui.Clear()

	// Print item title.
	msgs := makeDrawable(fmt.Sprintf("%c - %s", i.Hotkey(), i.Name()))
	PrintLong(msgs, 0)
	rows := len(msgs)

	// Print flavor text.
	msgs = makeDrawable(i.FlavorText())
	PrintLong(msgs, rows)
	rows += len(msgs)

	// Print flavor text.
	msgs = makeDrawable(StringEffects(i.Effects()))
	PrintLong(msgs, rows)
	rows += len(msgs)

	actionStr := dropAction
	if item.IsEquipable(i) && !creature.Hero.IsEquipped(i) {
		actionStr += " You can (e)quip this " + strings.ToLower(i.Cat()) + "."
	}
	if creature.Hero.IsEquipped(i) {
		actionStr += " You can (r)emove this " + strings.ToLower(i.Cat()) + "."
	}
	if item.IsUsable(i) {
		actionStr += " You can (u)se this " + strings.ToLower(i.Cat()) + "."
	}

	msgs = makeDrawable(actionStr)
	for y, m := range msgs {
		t := text.New(m, termbox.ColorCyan)
		ui.Print(t, ui.Inventory.XOffset, y+rows, ui.Inventory.Width)
	}

	termbox.Flush()

itemDetailLoop:
	for {
		switch detailsEvent := termbox.PollEvent(); detailsEvent.Type {
		case termbox.EventKey:
			if detailsEvent.Key == ui.CancelKey {
				break itemDetailLoop
			}

			switch string(detailsEvent.Ch) {
			case string(ui.DropItemKey):
				creature.Hero.DropItem(i.Hotkey(), a)
				return true
			case string(ui.EquipItemKey):
				NarrativeEquip(i.Hotkey())
				return true
			case string(ui.UseItemKey):
				NarrativeUse(i.Hotkey())
				return true
			case string(ui.UnEquipItemKey):
				NarrativeUnEquip(i.Hotkey())
				return true
			}
		}
	}
	return false
}

func makeDrawable(str string) []string {
	return strings.Split(stringsutil.WordWrap(str, ui.Inventory.Width), "\n")
}

func PrintLong(msg []string, yoffset int) {
	for y, m := range msg {
		t := text.New(m, termbox.ColorWhite)
		ui.Print(t, 0, y+yoffset, ui.Inventory.Width)
	}
	termbox.Flush()
}
