package gen

import (
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
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
