package ui

import (
	"github.com/karlek/worc/status"
)

const (
	// Game screen size.
	AreaScreenWidth  = 100
	AreaScreenHeight = 30

	// Status bar screen coordinates.
	statusX = 5
	statusY = 31

	// Status bar size.
	statusWidth  = 100 - statusX*2
	statusHeight = 7
)

func init() {
	// Init status menu.
	status.SetSize(statusWidth, statusHeight)
	status.SetLoc(statusX, statusY)
}
