package creature

import (
	"math"

	"github.com/karlek/reason/area"
	"github.com/karlek/reason/coord"
)

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
