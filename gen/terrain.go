package gen

import (
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
)

// Area is a debug function to generate terrain.
func Area(width, height int) area.Area {

	// Placeholder for terrain generation.
	var ms = []area.Tile{
		// fauna.Doodads["door (closed)"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["soil"],
		fauna.Doodads["bush"],
		fauna.Doodads["wall"],
		fauna.Doodads["water"],
		fauna.Doodads["water"],
		fauna.Doodads["water"],
		fauna.Doodads["water"],
		fauna.Doodads["water"],
		fauna.Doodads["water"],
	}

	a := area.New(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			a.Terrain[x][y] = ms[util.RandInt(0, len(ms))]
		}
	}
	return *a
}
