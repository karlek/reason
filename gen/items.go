package gen

import (
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
)

// Items is a debug function to add mobs to the map.
func Items(a *area.Area, num int) {
	var itemList = []item.Item{
		item.Items["Star-Eye Map"],
		item.Items["Ring of Orihalcon"],
	}
	for x := 0; x < a.Width; x++ {
		for y := 0; y < a.Height; y++ {
			if num <= 0 {
				return
			}

			if util.RandInt(0, 50) != 42 {
				continue
			}

			if !a.IsXYStackable(x, y) {
				continue
			}

			i := itemList[util.RandInt(0, len(itemList))]
			c := coord.Coord{x, y}

			i.NewX(x)
			i.NewY(y)

			if a.Items[c] == nil {
				a.Items[c] = new(area.Stack)
			}
			a.Items[c].Push(&i)
			num--
		}
	}
}
