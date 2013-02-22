package ui

import "github.com/karlek/worc/menu"

const (
	// Game screen size
	AreaScreenWidth  = 100
	AreaScreenHeight = 30

	// Status bar screen coordinates
	StatusX = 5
	StatusY = 31

	// Status bar size
	StatusWidth  = 100 - StatusX*2
	StatusHeight = 7
)

func init() {
	// Init status menu
	menu.SetStatusSize(StatusWidth, StatusHeight)
	menu.SetStatusLoc(StatusX, StatusY)
}

// type Keys struct {
// 	Open      termbox.Key
// 	Look      termbox.Key
// 	Quit      termbox.Key
// 	SaveQuit  termbox.Key
// 	MoveUp    termbox.Key
// 	MoveDown  termbox.Key
// 	MoveLeft  termbox.Key
// 	MoveRight termbox.Key
// }

// func LoadKeyBinds() {

// }
