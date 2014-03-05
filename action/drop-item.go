package action

import (
	"fmt"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"

	"github.com/nsf/termbox-go"
)

const (
	dropTitleFmt  = "Drop Item: %s (%s)"
	emptyInv      = "You aren't carrying anything."
	unableToEquip = "That item can't be equipped."
	unableToDrop  = "Couldn't drop item."
)

var (
	InvInfo = "currentWeight/maxPossibleWeight (usedSlots/totalSlots)"
)

func DropItem(a *area.Area) bool {
	// Show the inventory so the player knows which item to drop.
	title := fmt.Sprintf(dropTitleFmt, WeightStr, slotsString())
	isEmpty := categorizedInv(title)
	if isEmpty {
		return false
	}

	// Listen for user input to drop item.
	return dropInput(a)
}

func dropInput(a *area.Area) (actionTaken bool) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == ui.CancelKey {
				return false
			}
			if _, ok := creature.Hero.Inventory[ev.Ch]; ok {
				creature.Hero.DropItem(ev.Ch, a)
				return true
			}
		}
	}
}

func NarrativeEquip(pos rune) {
	i := creature.Hero.Equip(pos)
	if i == nil {
		status.Print(unableToEquip)
		return
	}

	var equipStr string
	if item.IsStackable(i) {
		equipStr += fmt.Sprintf("%d ", i.Count())
	}
	equipStr += i.Name()
	status.Printf("You equipped %s.", equipStr)
}

func NarrativeUse(pos rune) {
	creature.Hero.Use(creature.Hero.Inventory[pos])
}

func NarrativeUnEquip(pos rune) {
	creature.Hero.UnEquip(creature.Hero.Inventory[pos])
}

func InvText(i item.Itemer) string {
	invStr := ""
	if item.IsStackable(i) {
		invStr = fmt.Sprintf("%c - %d %s", i.Hotkey(), i.Count(), i.Name())
	} else {
		invStr = fmt.Sprintf("%c - %s", i.Hotkey(), i.Name())
	}

	if item.IsEquipable(i) {
		if creature.Hero.IsEquipped(i) {
			invStr += " (wielding)"
		}
	}
	return invStr
}

func InvAttr(i item.Itemer) termbox.Attribute {
	attr := RarityAttr(i)
	if item.IsEquipable(i) {
		if creature.Hero.IsEquipped(i) {
			attr = termbox.ColorGreen + termbox.AttrBold
		}
	}
	return attr
}

func RarityAttr(i item.Itemer) termbox.Attribute {
	var attr termbox.Attribute
	switch i.Rarity() {
	case item.Common:
		attr = termbox.ColorWhite
	case item.Magical:
		attr = termbox.ColorBlue + termbox.AttrBold
	case item.Artifact:
		attr = termbox.ColorWhite + termbox.AttrBold
	}
	return attr
}

func categorizedInv(title string) (isEmpty bool) {
	// Make screen black.
	ui.Clear()

	// If inventory is empty, warn the user.
	if len(creature.Hero.Inventory) == 0 {
		emptyInvMsg()

		// Show the text and wait for any input.
		termbox.Flush()
		termbox.PollEvent()
		return true
	}

	// Print inventory title.
	invTitle(title)

	// categories contains the item texts sorted in item categories.
	// Sort items into categories.
	var categories = map[string][]*ui.Text{}
	for _, pos := range item.Positions {
		if i, ok := creature.Hero.Inventory[pos]; ok {
			addToCategory(i, categories)
		}
	}

	// Print categories and the items in that category to screen.
	// Rows written to screen.
	rowOffset := 0
	for catStr, items := range categories {
		// Ignore empty categories
		if len(items) == 0 {
			continue
		}

		// Item category.
		printCategory(catStr, items, &rowOffset)
	}
	termbox.Flush()
	return false
}

func addToCategory(i item.Itemer, categories map[string][]*ui.Text) {
	t := ui.NewText(InvAttr(i), InvText(i))
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

func invTitle(title string) {
	t := ui.NewText(termbox.ColorWhite+termbox.AttrBold, title)
	ui.Print(t, 0, 0, ui.Whole.Width)
}

func emptyInvMsg() {
	t := ui.NewText(termbox.ColorWhite+termbox.AttrBold, emptyInv)
	ui.Print(t, 0, 0, ui.Whole.Width)
}

func printCategory(catStr string, items []*ui.Text, rowOffset *int) {
	cat := ui.NewText(termbox.ColorCyan+termbox.AttrBold, catStr)
	ui.Print(cat, 0, *rowOffset+ui.Inventory.YOffset, ui.Whole.Width)
	*rowOffset = *rowOffset + 1

	// Items in that category.
	for _, t := range items {
		ui.Print(t, ui.Inventory.XOffset, *rowOffset+ui.Inventory.YOffset, ui.Whole.Width)
		*rowOffset = *rowOffset + 1
	}

	// Empty line.
	*rowOffset = *rowOffset + 1
}
