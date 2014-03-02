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
	// ui.Area is the viewport size.

	camX := c.X() - ui.Area.Width/2
	camY := c.Y() - ui.Area.Height/2

	if c.X() < ui.Area.Width/2 {
		camX = 0
	}
	if c.Y() < ui.Area.Height/2 {
		camY = 0
	}
	if c.X() >= a.Width-ui.Area.Width/2 {
		camX = a.Width - ui.Area.Width
	}
	if c.Y() > a.Height-ui.Area.Height/2 {
		camY = a.Height - ui.Area.Height
	}

	return camX, camY
}
