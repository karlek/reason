package ui

import (
	"github.com/karlek/reason/area"
	"github.com/karlek/reason/coord"
)

// DrawFOV draws a field of view around a creature as well as the creatures
// memory of already explored areas.
func DrawFOV(c coord.Coord, cs map[coord.Coord]struct{}, a *area.Area) {
	// Clear screen.
	Clear()

	// Get viewport coordinate offset.
	camX, camY := camXY(c.X, c.Y, a)

	// Draw already explored areas.
	a.DrawExplored(Area, camX, camY)

	// Draw hero.
	a.Draw(c.X, c.Y, camX, camY, Area)

	// Visible coordinates of character.
	for p := range cs {
		// Set terrain as explored.
		a.Terrain[p.X][p.Y].IsExplored = true

		// TODO(_): refactor cam.
		a.Draw(p.X, p.Y, camX, camY, Area)
	}
}

// camXY returns the coordinate of offset for the viewport. Since the area can
// be larger than the viewport.
func camXY(x, y int, a *area.Area) (int, int) {
	cameraX, cameraY := camX(x, a), camY(y, a)
	if Area.Width > len(a.Terrain) {
		cameraX = 0
	}
	if Area.Height > len(a.Terrain[0]) {
		cameraY = 0
	}
	return cameraX, cameraY
}

func camX(x int, a *area.Area) int {
	// Area is the viewport size.
	cameraX := x - Area.Width/2

	if x < Area.Width/2 {
		cameraX = 0
	}
	if x >= a.Width-Area.Width/2 {
		cameraX = a.Width - Area.Width
	}
	return cameraX
}

func camY(y int, a *area.Area) int {
	cameraY := y - Area.Height/2
	if y < Area.Height/2 {
		cameraY = 0
	}
	if y > a.Height-Area.Height/2 {
		cameraY = a.Height - Area.Height
	}
	return cameraY
}
