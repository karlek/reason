// Package action implements actions for creatures.
package action

import (
	"github.com/karlek/reason/name"
	"github.com/karlek/reason/ui"

	"github.com/karlek/reason/ui/status"
	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/nsf/termbox-go"
)

// Look takes the units coordiantes and an area to observe.
func Look(a area.Area, x, y int) {
	x, y = x+ui.Area.XOffset, y+ui.Area.YOffset
	termbox.SetCursor(x, y)
	termbox.Flush()

lookLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Cursor movement.
			case ui.CancelKey:
				termbox.HideCursor()
				termbox.Flush()
				break lookLoop
			case ui.MoveUpKey:
				x, y = x, y-1
			case ui.MoveDownKey:
				x, y = x, y+1
			case ui.MoveLeftKey:
				x, y = x-1, y
			case ui.MoveRightKey:
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

			// Message is originally about terrain, but is then overwritten.
			msg := a.Terrain[x][y].(name.Namer).Name()
			if o, found := a.Objects[c]; found {
				if n, ok := o.(name.Namer); ok {
					msg = n.Name()
				}
			}
			if m, found := a.Monsters[c]; found {
				if n, ok := m.(name.Namer); ok {
					msg = n.Name()
				}
			}
			if s, found := a.Items[c]; found {
				if i := s.Peek(); i != nil {
					if n, ok := i.(name.Namer); ok {
						msg = n.Name()
					}
				}
			}
			status.Print("You see " + msg + ".")
		}
	}
}
