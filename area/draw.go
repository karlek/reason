package area

import (
	"github.com/karlek/reason/coord"
	"github.com/karlek/reason/draw"
	"github.com/karlek/reason/screen"

	"github.com/nsf/termbox-go"
)

func (a Area) Draw(x, y, cameraX, cameraY int, scr screen.Screen) {
	c := coord.Coord{X: x, Y: y}
	if m := a.Monsters[c]; m != nil {
		draw.DrawCell(x-cameraX, y-cameraY, m, scr)
		return
	}
	if o := a.Objects[c]; o != nil {
		draw.DrawCell(x-cameraX, y-cameraY, o, scr)
		return
	}
	if i := a.Items[c].Peek(); i != nil {
		draw.DrawCell(x-cameraX, y-cameraY, i, scr)
		return
	}
	draw.DrawCell(x-cameraX, y-cameraY, a.Terrain[x][y], scr)
}

func (a Area) DrawExplored(scr screen.Screen, cameraX, cameraY int) {
	for x := 0; x < a.Width; x++ {
		for y := 0; y < a.Height; y++ {
			if !a.Terrain[x][y].IsExplored {
				continue
			}
			a.drawMemory(x-cameraX, y-cameraY, cameraX, cameraY, scr)
		}
	}
}

func (a Area) drawMemory(x, y int, cameraX, cameraY int, scr screen.Screen) {
	c := coord.Coord{x + scr.XOffset, y + scr.YOffset}
	p := coord.Plane{scr.Width + scr.XOffset, scr.Height + scr.YOffset, scr.XOffset, scr.YOffset}
	if !p.Contains(c) {
		return
	}
	termbox.SetCell(x+scr.XOffset, y+scr.YOffset, a.Terrain[x+cameraX][y+cameraY].Graphic().Ch, termbox.ColorBlack+termbox.AttrBold, termbox.ColorDefault)
}
