package gen

import (
	"log"

	"github.com/karlek/reason/object"
	"github.com/karlek/reason/util"

	"github.com/karlek/reason/area"
	"github.com/karlek/reason/coord"
)

// Objects is a debug function to generate objects.
func Objects(a *area.Area, num int) {
	o := make(map[coord.Coord]area.DrawPather, a.Width*a.Height)

	log.Println(len(object.Objects))
	var objectList = make([]object.Object, len(object.Objects))
	var index int
	for _, o := range object.Objects {
		objectList[index] = o
		index++
	}
	for x := 0; x < a.Width; x++ {
		for y := 0; y < a.Height; y++ {
			if num <= 0 {
				return
			}

			if util.RandInt(0, 50) != 42 {
				continue
			}

			if !a.IsXYPathable(x, y) {
				continue
			}

			c := coord.Coord{X: x, Y: y}
			o[c] = objectList[util.RandInt(0, len(objectList))].New()
		}
	}
	a.Objects = o
}
