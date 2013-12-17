package gen

import (
	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
)

// Mobs is a debug function to add mobs to the map.
func Mobs(a *area.Area, num int) {
	var mobList = []beastiary.Creature{
		beastiary.Creatures["gobbal"],
		beastiary.Creatures["tofu"],
		beastiary.Creatures["arachnee"],
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

			g := mobList[util.RandInt(0, len(mobList))]
			c := coord.Coord{x, y}

			g.NewX(x)
			g.NewY(y)

			a.Monsters[c] = &g
			num--
		}
	}
	a.DrawMonsters()
}
