package action

import (
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

func CloseDoorNarrative(a *area.Area, x, y int) bool {
	status.Print("Close door - In which direction lies the door?")
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case ui.CancelKey:
			return false
		// Movement.
		case ui.MoveUpKey:
			x, y = x, y-1
		case ui.MoveDownKey:
			x, y = x, y+1
		case ui.MoveLeftKey:
			x, y = x-1, y
		case ui.MoveRightKey:
			x, y = x+1, y
		}
		// Prevent from moving outside the top
		if y < 0 {
			y = 0
		}

		// Prevent from moving outside the bottom
		if y == a.Height {
			y = a.Height - 1
		}

		// Prevent from moving outside the left
		if x < 0 {
			x = 0
		}

		// Prevent from moving outside the right
		if x == a.Width {
			x = a.Width - 1
		}

		if a.Terrain[x][y] == fauna.Doodads["door (open)"] {
			a.Terrain[x][y] = fauna.Doodads["door (closed)"]
			return true
		} else {
			status.Print("You can't close that.")
		}
	}
	return false
}
