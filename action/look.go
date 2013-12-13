// Package action implements actions for creatures.
package action

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/status"
	"github.com/nsf/termbox-go"
)

// Look takes the units coordiantes and an area to observe.
func Look(x, y int, a area.Area) {
	termbox.SetCursor(x, y)
	termbox.Flush()

lookLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Cursor movement.
			case termbox.KeyEsc:
				termbox.HideCursor()
				termbox.Flush()
				break lookLoop
			case termbox.KeyArrowUp:
				x, y = x, y-1
			case termbox.KeyArrowDown:
				x, y = x, y+1
			case termbox.KeyArrowLeft:
				x, y = x-1, y
			case termbox.KeyArrowRight:
				x, y = x+1, y
			default:
				continue
			}

			// Prevent from moving outside the top
			if y < 0 {
				y = 0
				continue
			}

			// Prevent from moving outside the bottom
			if y == a.Height {
				y = a.Height - 1
				continue
			}

			// Prevent from moving outside the left
			if x < 0 {
				x = 0
				continue
			}

			// Prevent from moving outside the right
			if x == a.Width {
				x = a.Width - 1
				continue
			}

			termbox.SetCursor(x, y)
			termbox.Flush()

			c := coord.Coord{x, y}
			var msg string
			if s, found := a.Objects[c]; found {
				// Object found message.
				if n, ok := s.Peek().(name.Namer); ok {
					msg = n.Name()
					if msg != "" {
						status.Print("You see " + msg + ".")
					}
					continue
				}
			}
			// If no object was found, print the terrains name instead.
			if msg == "" {
				msg = a.Terrain[x][y].(name.Namer).Name()
			}
			// Terrain found.
			status.Print("You see " + msg + ".")
		}
	}
}
