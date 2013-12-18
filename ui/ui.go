package ui

import (
	"fmt"
	"log"

	"github.com/karlek/progress/barcli"
	// "github.com/karlek/worc/area"
	"github.com/karlek/worc/screen"
	"github.com/nsf/termbox-go"
)

var (
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
		Width: 60,
	}
)

const (
	// Keybindings.
	// Action keys.
	LookKey          = 'l'
	OpenDoorKey      = 'o'
	CloseDoorKey     = 'c'
	QuitKey          = 'q'
	PickUpItemKey    = 'g'
	SaveAndQuitKey   = 'p'
	ShowInventoryKey = 'i'
	DropItemKey      = 'd'

	// Movement keys.
	MoveUpKey    = termbox.KeyArrowUp
	MoveDownKey  = termbox.KeyArrowDown
	MoveRightKey = termbox.KeyArrowRight
	MoveLeftKey  = termbox.KeyArrowLeft

	// General keys.
	CancelKey = termbox.KeyEsc
)

func print(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
	}
}

func PrintInventory(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
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
