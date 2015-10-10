package creature

import (
	"math"

	"github.com/karlek/reason/ui"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
)

// DrawFOV draws a field of view around a creature as well as the creatures
// memory of already explored areas.
func (c Creature) DrawFOV(a *area.Area) {
	// Clear screen.
	ui.Clear()

	// Get viewport coordinate offset.
	camX, camY := camXY(c, a)

	// Draw already explored areas.
	a.DrawExplored(ui.Area, camX, camY)

	// Draw hero.
	a.Draw(c.X(), c.Y(), camX, camY, ui.Area)

	// Visible coordinates of character.
	cs := c.FOV(a)
	for p := range cs {
		// Set terrain as explored.
		a.Terrain[p.X][p.Y].IsExplored = true

		// TODO(_): refactor cam.
		a.Draw(p.X, p.Y, camX, camY, ui.Area)
	}
}

func (c *Creature) FOV(a *area.Area) (cs map[coord.Coord]struct{}) {
	radius := c.Sight
	cs = make(map[coord.Coord]struct{})
	for x := c.X() - radius; x <= c.X()+radius; x++ {
		for y := c.Y() - radius; y <= c.Y()+radius; y++ {
			// Distance between creature x and y coordinates and sight radius.
			dx := float64(x - c.X())
			dy := float64(y - c.Y())

			// Distance between creature and sight radius.
			dist := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

			// Discriminate coordinates which are outside of the circle.
			if dist > float64(radius) {
				continue
			}

			// Ignore hero.
			for _, p := range get_line(c.X(), c.Y(), x, y)[1:] {
				if !a.ExistsXY(p.X, p.Y) {
					break
				}

				cs[p] = struct{}{}

				// Terrain that breaks line of sight.
				if !a.IsXYPathable(p.X, p.Y) {
					break
				}
			}
		}
	}
	return cs
}

func get_line(x1, y1, x2, y2 int) (points []coord.Coord) {
	points = make([]coord.Coord, 0)
	steep := math.Abs(float64(y2-y1)) > math.Abs(float64(x2-x1))
	if steep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}
	rev := false
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		rev = true
	}
	dx := x2 - x1
	dy := int(math.Abs(float64(y2 - y1)))
	err := dx / 2
	y := y1
	ystep := 0
	if y1 < y2 {
		ystep = 1
	} else {
		ystep = -1
	}
	for x := x1; x < x2+1; x++ {
		if steep {
			points = append(points, coord.Coord{X: y, Y: x})
		} else {
			points = append(points, coord.Coord{X: x, Y: y})
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
	if rev {
		reverse(points)
	}
	return points
}

func reverse(s []coord.Coord) []coord.Coord {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
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
