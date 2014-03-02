package action

import (
	"fmt"
	"strconv"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"

	"github.com/nsf/termbox-go"
)

const (
	DropTitleFmt  = "Drop Item: %s (%s)"
	EmptyInv      = "You aren't carrying anything."
	UnableToEquip = "That item can't be equipped."
	UnableToDrop  = "Couldn't drop item."
)

var (
	InvInfo = "currentWeight/maxPossibleWeight (usedSlots/totalSlots)"
)

func DropItem(a *area.Area) bool {
	categorizedInv(fmt.Sprintf(DropTitleFmt, WeightStr, slotsString()))
	if len(creature.Hero.Inventory) == 0 {
		return false
	}

dropItemLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == ui.CancelKey {
				break dropItemLoop
			}
			if _, ok := creature.Hero.Inventory[ev.Ch]; ok {
				creature.Hero.DropItem(ev.Ch, a)
				return true
			}
		}
	}
	return false
}

func NarrativeEquip(pos rune) {
	i := creature.Hero.Equip(pos)
	if i == nil {
		status.Print(UnableToEquip)
		return
	}

	s := "You equipped "
	if item.IsStackable(i) {
		s += strconv.Itoa(i.Count()) + " " + i.Name()
	} else {
		s += i.Name()
	}
	s += "."
	status.Print(s)
}

func NarrativeUse(pos rune) {
	creature.Hero.Use(creature.Hero.Inventory[pos])
}

func NarrativeUnEquip(pos rune) {
	creature.Hero.UnEquip(creature.Hero.Inventory[pos])
}

func InvText(i item.DrawItemer) string {
	s := ""
	if item.IsStackable(i) {
		s = string(i.Hotkey()) + " - " + strconv.Itoa(i.Count()) + " " + i.Name()
	} else {
		s = string(i.Hotkey()) + " - " + i.Name()
	}

	if item.IsEquipable(i) {
		if creature.Hero.IsEquipped(i) {
			s += " (wielding)"
		}
	}
	return s
}

func InvAttr(i item.DrawItemer) termbox.Attribute {
	attr := RarityAttr(i)
	if item.IsEquipable(i) {
		if creature.Hero.IsEquipped(i) {
			attr = termbox.ColorGreen + termbox.AttrBold
		}
	}
	return attr
}

func RarityAttr(i item.DrawItemer) termbox.Attribute {
	var Attr termbox.Attribute
	switch i.Rarity() {
	case item.Common:
		Attr = termbox.ColorWhite
	case item.Magical:
		Attr = termbox.ColorBlue + termbox.AttrBold
	case item.Artifact:
		Attr = termbox.ColorWhite + termbox.AttrBold
	}
	return Attr
}

func categorizedInv(title string) {
	// Make screen black.
	ui.Clear()

	// If inventory is empty, warn the user.
	if len(creature.Hero.Inventory) == 0 {
		t := ui.NewText(termbox.ColorWhite+termbox.AttrBold, EmptyInv)
		ui.Print(t, 0, 0, ui.Whole.Width)

		// Show the text and draw.
		termbox.Flush()
		termbox.PollEvent()
		return
	}

	// Print inventory title.
	t := ui.NewText(termbox.ColorWhite+termbox.AttrBold, title)
	ui.Print(t, 0, 0, ui.Whole.Width)

	// categories contains the item texts sorted in item categories.
	var categories = map[string][]*ui.Text{}

	// Sort items into categories.
	for _, pos := range item.Positions {
		if i, ok := creature.Hero.Inventory[pos]; ok {
			t = ui.NewText(InvAttr(i), InvText(i))
			switch i.(type) {
			case *item.Weapon:
				categories["Weapons"] = append(categories["Weapons"], t)
			case *item.Tool:
				categories["Tools"] = append(categories["Tools"], t)
			case *item.Ring:
				categories["Rings"] = append(categories["Rings"], t)
			case *item.Potion:
				categories["Potions"] = append(categories["Potions"], t)
			default:
				categories["Unknown"] = append(categories["Unknown"], t)
			}
		}
	}

	// Rows written to screen.
	rowOffset := 0

	// Print categories and the items in that category to screen.
	for catStr, items := range categories {
		if len(items) == 0 {
			continue
		}

		// Item category.
		cat := ui.NewText(termbox.ColorCyan+termbox.AttrBold, catStr)
		ui.Print(cat, 0, rowOffset+ui.Inventory.YOffset, ui.Whole.Width)
		rowOffset++

		// Items in that category.
		for _, t := range items {
			ui.Print(t, ui.Inventory.XOffset, rowOffset+ui.Inventory.YOffset, ui.Whole.Width)
			rowOffset++
		}

		// Empty line.
		rowOffset++
	}
	termbox.Flush()
}
