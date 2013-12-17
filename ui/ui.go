package ui

import (
	"fmt"

	"github.com/karlek/reason/beastiary"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

const (
	// Whole game screen.
	WholeScreenWidth  = 127
	WholeScreenHeight = 44

	// Area screen size.
	AreaScreenWidth  = 75
	AreaScreenHeight = 36

	// Character info.
	CharacterInfoX      = AreaScreenWidth + 2
	CharacterInfoY      = 0
	CharacterInfoWidth  = 25
	CharacterInfoHeight = AreaScreenHeight

	// Status bar screen coordinates.
	statusX = 5
	statusY = AreaScreenHeight + 1

	// Status bar size.
	statusWidth  = 100 - statusX*2
	statusHeight = 7

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

func init() {
	// Init status menu.
	status.SetSize(statusWidth, statusHeight)
	status.SetLoc(statusX, statusY)
}

func print(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
	}
	termbox.Flush()
}

func PrintInventory(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
	}
	termbox.Flush()
}

// UpdateHp updates the hero health bar.
func UpdateHp(hero beastiary.Creature) {
	hpMsg := fmt.Sprintf("%d/%d", hero.Hp, hero.MaxHp)
	print("Health: ", CharacterInfoX, CharacterInfoY+1, CharacterInfoWidth, termbox.ColorWhite, termbox.ColorDefault)
	print(hpMsg, CharacterInfoX+8, CharacterInfoY+1, CharacterInfoWidth, termbox.ColorRed, termbox.ColorDefault)
}

// Update draws all user interfaces to the screen.
func Update(hero beastiary.Creature) {
	UpdateHp(hero)
}

func AreaScreenRedraw(a area.Area, hero beastiary.Creature) {
	a.Draw()
	Update(hero)
	status.Update()
}
