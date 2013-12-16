package ui

import (
	"fmt"

	"github.com/karlek/reason/beastiary"

	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

const (
	// Game screen size.
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

// UpdateHp updates the hero health bar.
func UpdateHp(hero beastiary.Creature) {
	hpMsg := fmt.Sprintf("%d/%d", hero.Hp, hero.MaxHp)
	print("Health: ", CharacterInfoX, CharacterInfoY+1, CharacterInfoWidth, termbox.ColorWhite, termbox.ColorDefault)
	print(hpMsg, CharacterInfoX+8, CharacterInfoY+1, CharacterInfoWidth, termbox.ColorRed, termbox.ColorDefault)
}

// Init draws all user interfaces to the screen.
func Init(hero beastiary.Creature) {
	UpdateHp(hero)
}
