package ui

import (
	"fmt"
	"log"

	"github.com/karlek/reason/ui/text"

	"github.com/karlek/progress/barcli"
	"github.com/karlek/worc/screen"
	"github.com/nsf/termbox-go"
)

var (
	Terminal = screen.Screen{}

	// Main menu screen size.
	Main = screen.Screen{
		Width:   35,
		Height:  20,
		YOffset: 50,
	}

	// Area screen size.
	// Area = screen.Screen{
	// 	Width:  35,
	// 	Height: 20,
	// }
	Area = screen.Screen{
		Width:  50,
		Height: 30,
	}
	CharacterInfo = screen.Screen{}

	// CharacterInfo = screen.Screen{
	// 	Width:   25,
	// 	Height:  Area.Height,
	// 	YOffset: 0,
	// 	XOffset: 1,
	// }

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

func SetTerminal() {
	Terminal.Width, Terminal.Height = termbox.Size()
	Area = screen.Screen{
		Width:  Area.Width,
		Height: Area.Height,
		// YOffset: Terminal.Height/2 - Area.Height + Area.Height/4,
		// XOffset: Terminal.Width/2 - Area.Width + Area.Width/4,
		YOffset: 2,
		XOffset: Terminal.Width/2 - Area.Width/2,
	}
	CharacterInfo = screen.Screen{
		Width:   25,
		Height:  Area.Height,
		YOffset: 0,
		XOffset: Area.XOffset,
	}
	MonsterInfo = screen.Screen{
		Width:   20,
		Height:  Area.Height,
		YOffset: Area.YOffset,
		XOffset: Area.XOffset + Area.Width + 2,
	}
	Message = screen.Screen{
		Width:   Whole.Width,
		Height:  5,
		YOffset: Area.Height + Area.YOffset + 1,
		XOffset: Area.XOffset + 1,
	}
}

func Clear() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
}

// // clearLine clears the line from old characters.
// func ClearLine(x, y, length, width int) {
// 	for i := length; i < width; i++ {
// 		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
// 	}
// }

func ClearLine(line int) {
	for x := 0; x < Terminal.Width; x++ {
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

func Print(t *text.Text, x, y, width int) {
	print(t.Text, x, y, width, t.Attr, termbox.ColorDefault)
}

// Hp updates the hero health bar.
// Health: 17/17 =====================
func Hp(curHp, maxHp int) {
	xOffset := 0

	t := text.New("Health: ", termbox.ColorWhite)
	Print(
		t,
		CharacterInfo.XOffset,
		CharacterInfo.YOffset,
		CharacterInfo.Width,
	)
	// len("Health: ")
	xOffset += 8

	t.Text = fmt.Sprintf("%d/%d", curHp, maxHp)
	t.Attr = termbox.ColorRed
	Print(
		t,
		CharacterInfo.XOffset+xOffset,
		CharacterInfo.YOffset,
		CharacterInfo.Width,
	)
	xOffset += len(t.Text) + 1

	bar, err := barcli.New(maxHp)
	if err != nil {
		log.Println(err)
	}
	bar.IncN(curHp)

	filled, unfilled, err := bar.StringSize(20)
	if err != nil {
		log.Println(err)
	}

	t.Text = filled
	t.Attr = termbox.ColorGreen + termbox.AttrBold
	Print(
		t,
		CharacterInfo.XOffset+xOffset,
		CharacterInfo.YOffset,
		CharacterInfo.Width,
	)
	xOffset += len(filled)

	t.Text = unfilled
	t.Attr = termbox.ColorBlack + termbox.AttrBold
	Print(
		t,
		CharacterInfo.XOffset+xOffset,
		CharacterInfo.YOffset,
		CharacterInfo.Width,
	)
}
