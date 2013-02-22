package beastiary

import "github.com/karlek/worc/creature"
import "github.com/nsf/termbox-go"

var Hero = creature.Creature{
	Graphic: termbox.Cell{
		Ch: '@',
		Fg: termbox.ColorWhite + termbox.AttrBold,
	},
	Name: "Hero",
	X:    50,
	Y:    15,
	Stats: map[string]int{
		"Health":   10,
		"Strength": 3,
	},
}

var Gobbal = creature.Creature{
	Graphic: termbox.Cell{
		Ch: 'G',
		Fg: termbox.ColorWhite + termbox.AttrBold,
	},
	Name: "Gobbal",
	Stats: map[string]int{
		"Health":   10,
		"Strength": 3,
	},
}
