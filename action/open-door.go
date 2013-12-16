// Package action implements actions for creatures.
package action

import (
	"github.com/karlek/reason/fauna"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

// func OpenDoor(a *area.Area, x, y int) bool {
// }

// ///
// func openDoor(a *area.Area, x, y int) bool {

// }

func OpenDoorNarrative(a *area.Area, x, y int) bool {
	status.Print("Toggle door - In which direction lies the door?")
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyEsc:
			return false
		// Movement
		case termbox.KeyArrowUp:
			x, y = x, y-1
		case termbox.KeyArrowDown:
			x, y = x, y+1
		case termbox.KeyArrowLeft:
			x, y = x-1, y
		case termbox.KeyArrowRight:
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

		if a.Terrain[x][y] == fauna.Doodads["door (closed)"] {
			a.Terrain[x][y] = fauna.Doodads["door (open)"]
			a.ReDraw(x, y)
			return true
		} else if a.Terrain[x][y] == fauna.Doodads["door (open)"] {
			a.Terrain[x][y] = fauna.Doodads["door (closed)"]
			a.ReDraw(x, y)
			return true
		} else {
			status.Print("You can't open / close that.")
		}
	}
	return false
}

func WalkedIntoDoor(a *area.Area, x, y int) bool {
	status.Print("Do you want to open door? [Y/n]")
	wantToOpenDoor := YesOrNo()
	if wantToOpenDoor {
		a.Terrain[x][y] = fauna.Doodads["door (open)"]
		a.ReDraw(x, y)
		return true
	}
	return false
}

// YesOrNo forces the player to answer either y or n.
// Esc is false and enter is true.
func YesOrNo() bool {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'n':
				return false
			case 'y':
				return true
			}
			switch ev.Key {
			case termbox.KeyEsc:
				return false
			case termbox.KeyEnter:
				return true
			}
		}
	}
}
