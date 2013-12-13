package gen

import (
	"math/rand"
	"time"

	"github.com/karlek/reason/fauna"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/screen"
)

// Area is a debug function to generate terrain.
func Area(scr screen.Screen, width, height int) area.Area {
	// Placeholder for terrain generation.
	var ms = []area.Stackable{
		&fauna.Soil,
		&fauna.Soil,
		&fauna.Soil,
		&fauna.Soil,
		&fauna.Soil,
		&fauna.Soil,
		&fauna.Wall,
	}

	a := area.Area{
		Terrain: make([][]area.Stackable, width),
		Objects: make(map[coord.Coord]*area.Stack),
		Width:   width,
		Height:  height,
		Screen:  scr,
	}

	for x := 0; x < width; x++ {
		a.Terrain[x] = make([]area.Stackable, height)
		for y := 0; y < height; y++ {
			a.Terrain[x][y] = ms[randInt(0, len(ms))]
		}
	}
	return a
}

// randInt is used by the debug function GenArea.
func randInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
