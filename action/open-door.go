// Package action implements actions for creatures.
package action

import (
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/util"

	"github.com/karlek/reason/ui/status"
	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

func direction(ev termbox.Event, x, y *int, a *area.Area) bool {
	switch ev.Key {
	case ui.CancelKey:
		return true
	// Movement
	case ui.MoveUpKey:
		*x, *y = *x, *y-1
	case ui.MoveDownKey:
		*x, *y = *x, *y+1
	case ui.MoveLeftKey:
		*x, *y = *x-1, *y
	case ui.MoveRightKey:
		*x, *y = *x+1, *y
	}
	// Prevent from moving outside the top
	if *y < 0 {
		*y = 0
	}

	// Prevent from moving outside the bottom
	if *y == a.Height {
		*y = a.Height - 1
	}

	// Prevent from moving outside the left
	if *x < 0 {
		*x = 0
	}

	// Prevent from moving outside the right
	if *x == a.Width {
		*x = a.Width - 1
	}
	return false
}

func OpenDoorNarrative(a *area.Area, x, y int) bool {
	status.Print("Open door - In which direction lies the door?")
	termbox.Flush()

	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		cancel := direction(ev, &x, &y, a)
		if cancel {
			return false
		}
		if a.Terrain[x][y] == fauna.Doodads["door (closed)"] {
			a.Terrain[x][y] = fauna.Doodads["door (open)"]
			return true
		}
	}
	status.Printf("You can't open that.")
	status.Printf("%T", a.Terrain[x][y])

	return false
}

func WalkedIntoDoor(a *area.Area, x, y int) bool {
	status.Print("Do you want to open door? [Y/n]")
	termbox.Flush()

	// user wants to open door?
	if util.YesOrNo() {
		a.Terrain[x][y] = fauna.Doodads["door (open)"]
		return true
	}

	return false
}
