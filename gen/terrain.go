package gen

import (
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/object"
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/util"

	"github.com/karlek/reason/area"
	"github.com/karlek/reason/coord"
)

// Area is a debug function to generate terrain.
func Area(width, height int) area.Area {

	// Placeholder for terrain generation.
	var ms = []terrain.Terrain{
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["soil"],
		terrain.Fauna["bush"],
		terrain.Fauna["wall"],
		terrain.Fauna["water"],
		terrain.Fauna["water"],
		terrain.Fauna["water"],
		terrain.Fauna["water"],
		terrain.Fauna["water"],
		terrain.Fauna["water"],
	}

	a := area.New(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			a.Terrain[x][y] = terrain.New(ms[util.RandInt(0, len(ms))])
		}
	}
	return *a
}

func translate(r rune) *terrain.Terrain {
	switch r {
	case '#':
		return terrain.New(terrain.Fauna["wall"])
	case '&':
		return terrain.New(terrain.Fauna["bush"])
	case '~':
		return terrain.New(terrain.Fauna["water"])
	case '.':
		return terrain.New(terrain.Fauna["soil"])
	default:
		return terrain.New(terrain.Fauna["soil"])
	}
}

func AreaPrim(width, height int) area.Area {
	raw := []string{
		`######################################################`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#.......................#######......................#`,
		`#.......................#.....#......................#`,
		`#.......................#....&#......................#`,
		`#.......................#.&.~.####...................#`,
		`#.......................#.~~~~#).#...................#`,
		`#.......................#~~~~~##+#...................#`,
		`#.......................##&..&#......................#`,
		`#...........................&.#......................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`#....................................................#`,
		`######################################################`,
	}

	a := area.New(len(raw[0]), len(raw))
	for y, row := range raw {
		for x, cell := range row {
			if cell == ')' {
				c := coord.Coord{X: x, Y: y}
				if stk, ok := a.Items[c]; ok {
					stk.Push(item.Items["Iron Sword"])
				} else {
					a.Items[c] = new(area.Stack)
					a.Items[c].Push(item.Items["Iron Sword"])
				}
			}
			if cell == '+' {
				a.Objects[coord.Coord{X: x, Y: y}] = object.Objects["door (closed)"].New()
			}
			a.Terrain[x][y] = translate(cell)
		}
	}
	return *a
}
