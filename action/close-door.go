package action

import (
	// "github.com/karlek/reason/fauna"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

func CloseDoorNarrative(a *area.Area, x, y int) bool {
	status.Print("Close door - In which direction lies the door?")
	termbox.Flush()

	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		cancel := direction(ev, &x, &y, a)
		if cancel {
			return false
		}
		return true
		// if a.Terrain[x][y] == fauna.Doodads["door (open)"] {
		// 	a.Terrain[x][y] = fauna.Doodads["door (closed)"]
		// 	return true
		// }
	}
	status.Print("You can't close that.")
	return false
}
