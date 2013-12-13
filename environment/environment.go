package environment

import "github.com/karlek/worc/terrain"
import "github.com/nsf/termbox-go"

var Tree = terrain.Terrain{
	Graphic: termbox.Cell{
		Ch: 'T',
		Fg: termbox.ColorGreen + termbox.AttrBold,
	},
	Name:     "a tree",
	Moveable: false,
}

var Wall = terrain.Terrain{
	Graphic: termbox.Cell{
		Ch: '#',
		Fg: termbox.ColorWhite + termbox.AttrBold,
	},
	Name:     "a wall",
	Moveable: false,
}

var Soil = terrain.Terrain{
	Graphic: termbox.Cell{
		Ch: '.',
		Fg: termbox.ColorYellow,
	},
	Name:     "soil",
	Moveable: true,
}

var ClosedDoor = terrain.Terrain{
	Graphic: termbox.Cell{
		Ch: '+',
		Fg: termbox.ColorYellow,
	},
	Name:     "a closed door",
	Moveable: false,
}

var OpenDoor = terrain.Terrain{
	Graphic: termbox.Cell{
		Ch: '-',
		Fg: termbox.ColorYellow,
	},
	Name:     "an open door",
	Moveable: true,
}
