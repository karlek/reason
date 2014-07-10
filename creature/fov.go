package creature

import (
	"math"

	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"
)

// DrawFOV draws a field of view around a creature as well as the creatures
// memory of already explored areas.
func (c Creature) DrawFOV(a *area.Area) {
	ui.Clear()

	// Get viewport coordinate offset.
	camX, camY := camXY(c, a)

	// Draw already explored areas.
	a.DrawExplored(ui.Area, camX, camY)

	// Inclusive hero's square so it's from hero's eyes.
	radius := c.Sight
	for x := c.X() - radius; x <= c.X()+radius; x++ {
		for y := c.Y() - radius; y <= c.Y()+radius; y++ {
			// Discriminate coordinates which are out of bounds.
			if !a.ExistsXY(x, y) {
				continue
			}

			// Distance between creature x and y coordinates and sight radius.
			dx := float64(x - c.X())
			dy := float64(y - c.Y())

			// Distance between creature and sight radius.
			dist := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

			// Discriminate coordinates which are outside of the circle.
			if dist > float64(radius) {
				continue
			}

			// Set terrain as explored.
			a.Terrain[x][y].IsExplored = true

			// TODO(_): refactor cam.
			a.Draw(x, y, camX, camY, ui.Area)
		}
	}
}

// camXY returns the coordinate of offset for the viewport. Since the area can
// be larger than the viewport.
func camXY(c Creature, a *area.Area) (int, int) {
	cameraX, cameraY := camX(c, a), camY(c, a)
	if ui.Area.Width > len(a.Terrain) {
		cameraX = 0
	}
	if ui.Area.Height > len(a.Terrain[0]) {
		cameraY = 0
	}
	return cameraX, cameraY
}

func camX(c Creature, a *area.Area) int {
	// ui.Area is the viewport size.
	cameraX := c.X() - ui.Area.Width/2

	if c.X() < ui.Area.Width/2 {
		cameraX = 0
	}
	if c.X() >= a.Width-ui.Area.Width/2 {
		cameraX = a.Width - ui.Area.Width
	}
	return cameraX
}

func camY(c Creature, a *area.Area) int {
	cameraY := c.Y() - ui.Area.Height/2
	if c.Y() < ui.Area.Height/2 {
		cameraY = 0
	}
	if c.Y() > a.Height-ui.Area.Height/2 {
		cameraY = a.Height - ui.Area.Height
	}
	return cameraY
}
