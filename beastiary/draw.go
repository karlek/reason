package beastiary

import (
	"math"

	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"
	// "github.com/karlek/worc/screen"
	"github.com/nsf/termbox-go"
)

func (c Creature) DrawFOV(a *area.Area) {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	cameraX, cameraY := CameraXY(c, a)
	a.DrawExplored(ui.Area, cameraX, cameraY)

	radius := 9 // Inclusive hero's square, so actually 8 from hero's "eyes".
	for x := c.X() - radius; x <= c.X()+radius; x++ {
		for y := c.Y() - radius; y <= c.Y()+radius; y++ {
			if !a.ExistsXY(x, y) {
				continue
			}

			// Distance between creature x and y coordinates and sight radius.
			dx := float64(x - c.X())
			dy := float64(y - c.Y())

			// Distance between creature and sight radius.
			dist := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

			if dist > float64(radius) {
				continue
			}
			// / workaround for pointer reciver on area.Tile problem.
			tile := a.Terrain[x][y]
			d, ok := tile.(fauna.Doodad)
			if !ok {
				continue
			}
			d.Explored = true
			a.Terrain[x][y] = d

			a.Draw(x, y, cameraX, cameraY, ui.Area)
		}
	}
}

func CameraXY(c Creature, a *area.Area) (int, int) {
	cameraX := c.X() - (ui.Area.Width / 2)
	cameraY := c.Y() - (ui.Area.Height / 2)

	if c.X() < (ui.Area.Width / 2) {
		cameraX = 0
	}
	if c.Y() < (ui.Area.Height / 2) {
		cameraY = 0
	}
	if c.X() >= a.Width-(ui.Area.Width/2) {
		cameraX = a.Width - ui.Area.Width
	}
	if c.Y() > a.Height-(ui.Area.Height/2) {
		cameraY = a.Height - ui.Area.Height
	}

	return cameraX, cameraY
}
