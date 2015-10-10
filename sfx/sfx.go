package sfx

import (
	"time"

	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/util"

	"github.com/nsf/termbox-go"
)

func glitch(amount int) {
	for y := 0; y < ui.Terminal.Height; y++ {
		for x := 0; x < ui.Terminal.Width; x++ {
			if util.RandInt(0, amount) != 1 {
				continue
			}
			c1 := util.RandInt(0, len(colors))
			c2 := util.RandInt(0, len(colors))
			a := util.RandInt(0, len(atts))
			termbox.SetCell(x, y, rune(util.RandInt(0, 255)), colors[c2]+atts[a], colors[c1]+atts[a])
		}
	}
	termbox.Flush()
}

func Glitch() {
	glitch(16)
	time.Sleep(time.Millisecond * 120)
	glitch(7)
	time.Sleep(time.Millisecond * 60)
	glitch(2)
}

var atts = []termbox.Attribute{
	termbox.AttrBold,
	termbox.AttrUnderline,
	termbox.AttrReverse,
}

var colors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorBlue,
	termbox.ColorCyan,
	termbox.ColorDefault,
	termbox.ColorGreen,
	termbox.ColorMagenta,
	termbox.ColorRed,
	termbox.ColorWhite,
	termbox.ColorYellow,
}
