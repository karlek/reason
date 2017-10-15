package gen

import (
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/util"

	"github.com/karlek/reason/area"
	"github.com/karlek/reason/coord"
)

// Items is a debug function to add mobs to the map.
func Items(a *area.Area, num int) {
	var itemList = make([]item.DrawItemer, len(item.Items))
	var index int
	for _, i := range item.Items {
		itemList[index] = i
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

			c := coord.Coord{x, y}
			i := item.New(itemList[util.RandInt(0, len(itemList))])
			if i == nil {
				continue
			}
			i.SetX(x)
			i.SetY(y)

			if a.Items[c] == nil {
				a.Items[c] = new(area.Stack)
			}
			a.Items[c].Push(i)
			num--
		}
	}
}
