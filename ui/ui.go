package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/karlek/progress/barcli"
	"github.com/karlek/worc/screen"
	"github.com/nsf/termbox-go"
)

var (
	// Main menu screen size.
	Main = screen.Screen{
		Width:   35,
		Height:  20,
		YOffset: 50,
	}

	// Area screen size.
	Area = screen.Screen{
		Width:   35,
		Height:  20,
		YOffset: 2,
	}

	CharacterInfo = screen.Screen{
		Width:   25,
		Height:  Area.Height,
		YOffset: 0,
		XOffset: 1,
	}

	MonsterInfo = screen.Screen{
		Width:   20,
		Height:  Area.Height,
		YOffset: Area.YOffset,
		XOffset: Area.Width + 2,
	}

	Message = screen.Screen{
		Width:   Whole.Width,
		Height:  5,
		YOffset: Area.Height + Area.YOffset + 1,
		XOffset: 1,
	}

	Whole = screen.Screen{
		Width:  34 + 25 + 2,
		Height: 44,
	}

	Inventory = screen.Screen{
		Width:   105,
		YOffset: 2,
		XOffset: 1,
	}
)

const (
	// Keybindings.
	// Action keys.
	CloseDoorKey     = 'c'
	DropItemKey      = 'd'
	EquipItemKey     = 'e'
	UnEquipItemKey   = 'r'
	UseItemKey       = 'u'
	LookKey          = 'l'
	OpenDoorKey      = 'o'
	PickUpItemKey    = 'g'
	QuitKey          = 'q'
	SaveAndQuitKey   = 'p'
	ShowInventoryKey = 'i'

	// Movement keys.
	MoveUpKey    = termbox.KeyArrowUp
	MoveDownKey  = termbox.KeyArrowDown
	MoveRightKey = termbox.KeyArrowRight
	MoveLeftKey  = termbox.KeyArrowLeft

	// General keys.
	CancelKey = termbox.KeyEsc
)

type Text struct {
	Attr termbox.Attribute
	Text string
}

func NewText(attr termbox.Attribute, str string) *Text {
	return &Text{Attr: attr, Text: str}
}

func Clear() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

func ClearLine(line int) {
	/// Make relative to window size.
	for x := 0; x < 170; x++ {
		termbox.SetCell(x, line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}

func ClearLineOffset(line, x int) {
	/// Make relative to window size.
	for ; x < 170; x++ {
		termbox.SetCell(x, line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}

func print(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
	}
}

func Print(t *Text, x, y, width int) {
	// Clears the line from old characters.
	for i := len(t.Text); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range t.Text {
		termbox.SetCell(x+charOffset, y, char, t.Attr, termbox.ColorDefault)
	}
}

// UpdateHp updates the hero health bar.
func UpdateHp(curHp, maxHp int) {
	hpMsg := fmt.Sprintf("%d/%d", curHp, maxHp)

	xOffset := 0
	print("Health: ", CharacterInfo.XOffset, CharacterInfo.YOffset, CharacterInfo.Width, termbox.ColorWhite, termbox.ColorDefault)
	xOffset += 8
	print(hpMsg, CharacterInfo.XOffset+xOffset, CharacterInfo.YOffset, CharacterInfo.Width, termbox.ColorRed, termbox.ColorDefault)
	xOffset += len(hpMsg) + 1

	bar, err := barcli.New(maxHp)
	if err != nil {
		log.Println(err)
	}
	bar.IncN(curHp)
	filled, unfilled, err := bar.StringSize(20)
	if err != nil {
		log.Println(err)
	}

	print(filled, CharacterInfo.XOffset+xOffset, CharacterInfo.YOffset, CharacterInfo.Width, termbox.ColorGreen+termbox.AttrBold, termbox.ColorDefault)
	xOffset += len(filled)
	print(unfilled, CharacterInfo.XOffset+xOffset, CharacterInfo.YOffset, CharacterInfo.Width, termbox.ColorBlack+termbox.AttrBold, termbox.ColorDefault)
}

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

func UpdateMonsterInfo(info []MonstInfo) {
	for yOffset, monst := range info[:] {
		t := NewText(monst.Graphics.Fg, string(monst.Graphics.Ch))
		Print(t, MonsterInfo.XOffset, MonsterInfo.YOffset+yOffset, MonsterInfo.Width)

		t = NewText(monst.Color(), "█")
		Print(t, MonsterInfo.XOffset+2, MonsterInfo.YOffset+yOffset, MonsterInfo.Width)

		t = NewText(termbox.ColorWhite, strings.Title(monst.Name))
		Print(t, MonsterInfo.XOffset+4, MonsterInfo.YOffset+yOffset, MonsterInfo.Width)
	}
}
