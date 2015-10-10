package ui

import (
	"strings"

	"github.com/karlek/reason/ui/text"

	"github.com/nsf/termbox-go"
)

type MonstInfo struct {
	Name     string
	HpLevel  int
	Graphics termbox.Cell
}

func (mi MonstInfo) Color() termbox.Attribute {
	switch mi.HpLevel {
	case 1:
		return termbox.ColorGreen + termbox.AttrBold
	case 2:
		return termbox.ColorYellow + termbox.AttrBold
	case 3:
		return termbox.ColorMagenta + termbox.AttrBold
	case 4:
		return termbox.ColorRed
	}
	return termbox.ColorBlack
}

func Monsters(info []MonstInfo) {
	// G █ Name
	//
	// G is creatures graphic.
	// █ is the color representation of the health of the monster.
	// Name is the name of the creature.
	for y, monst := range info {
		t := text.New(string(monst.Graphics.Ch), monst.Graphics.Fg)
		Print(
			t,
			monsterInfo.XOffset,
			monsterInfo.YOffset+y,
			monsterInfo.Width,
		)

		t = text.New("█", monst.Color())
		Print(
			t,
			monsterInfo.XOffset+2,
			monsterInfo.YOffset+y,
			monsterInfo.Width,
		)

		t = text.New(strings.Title(monst.Name), termbox.ColorWhite)
		Print(
			t,
			monsterInfo.XOffset+4,
			monsterInfo.YOffset+y,
			monsterInfo.Width,
		)
	}
}
