package intro

import (
	"log"
	"os"
	"unicode"

	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/util"

	"github.com/nsf/termbox-go"
)

const (
	journeyStr = "Do you really need a reason to begin your journey?"
	nameStr    = "What's your name?"

	// Credits to Fraktur @ http://patorjk.com/software/taag
	logo = `` +
		`     ..      ...                               .x+=:.                             ` +
		"  :~\"8888x :\"%888x                            z`    ^%                            " +
		` 8    8888Xf  8888>                              .   <k        u.      u.    u.   ` +
		`X88x. ?8888k  8888X       .u          u        .@8Ned8"  ...ue888b   x@88k u@88c. ` +
		`'8888L'8888X  '%88X    ud8888.     us888u.   .@^%8888"   888R Y888r ^"8888""8888" ` +
		" \"888X 8888X:xnHH(`` :888'8888. .@88 \"8888\" x88:  `)8b.  888R I888>   8888  888R  " +
		`   ?8~ 8888X X8888   d888 '88%" 9888  9888  8888N=*8888  888R I888>   8888  888R  ` +
		" -~`   8888> X8888   8888.+\"    9888  9888   %8\"    R88  888R I888>   8888  888R  " +
		` :H8x  8888  X8888   8888L      9888  9888    @8Wou 9%  u8888cJ888    8888  888R  ` +
		" 8888> 888~  X8888   '8888c. .+ 9888  9888  .888888P`    \"*888*P\"    \"*88*\" 8888\" " +
		" 48\"` '8*~   `8888!`  \"88888%   \"888*\"\"888\" `   ^\"F        'Y\"         \"\"   'Y\"   " +
		"  ^-==\"\"      `\"\"       \"YP'     ^Y\"   ^Y'                                        "
)

var (
	name string
)

func Intro() string {
	printJourney()
	printLogo()
	printEnterSign()
	waitForInput()
	removeEnterSign()
	printNameQuestion()
	return askName()
}

func printLogo() {
	width, height := termbox.Size()
	logoH := 12
	logoW := 82

	for x := 0; x < logoW; x++ {
		for y := 0; y < logoH; y++ {
			t := ui.NewText(termbox.ColorRed, string(logo[y*logoW+x]))
			ui.Print(t, width/2-logoW/2+x, height/2-15+y, 0)
		}
	}
	termbox.Flush()
}

func printEnterSign() {
	width, height := termbox.Size()

	msgPos := width/2 + len(journeyStr)/2
	t := ui.NewText(termbox.AttrBold, "↵")
	ui.Print(t, msgPos+1, height/2, 0)
	termbox.Flush()
}

func removeEnterSign() {
	width, height := termbox.Size()
	msgPos := width/2 + len(journeyStr)/2

	ui.ClearLineOffset(height/2, msgPos)
}

func printJourney() {
	width, height := termbox.Size()
	msgPos := width/2 - len(journeyStr)/2

	ui.ClearLineOffset(height/2, msgPos)

	t := ui.NewText(termbox.AttrBold, journeyStr)
	ui.Print(t, msgPos, height/2, 0)
}

func printNameQuestion() {
	width, height := termbox.Size()

	offsetY := height/2 + 2
	msgPos := width/2 - len(journeyStr)/2

	t := ui.NewText(termbox.AttrBold, nameStr)
	ui.Print(t, msgPos, offsetY, 0)
}

func printName() {
	width, height := termbox.Size()

	offsetX := width/2 - len(nameStr)/2
	offsetY := height/2 + 2
	ui.ClearLineOffset(offsetY, offsetX)

	t := ui.NewText(termbox.AttrBold+termbox.ColorBlack, name)
	ui.Print(t, offsetX+1, offsetY, 0)
}

func waitForInput() {
	ev := termbox.PollEvent()
	if ev.Key == ui.CancelKey {
		termbox.Close()
		os.Exit(0)
	}
}

func backspace() {
	// Remove the last character from name.
	if len(name) > 0 {
		name = name[:len(name)-1]
	}
	printNameQuestion()
	printName()
}

func addToName(r rune) {
	// Ignore non-printable characters or if the name is to long.
	if !unicode.IsPrint(r) || len(name) > 30 {
		log.Printf("%v\n", r)
		return
	}
	name += string(r)
	printName()
}

func askName() string {
	for {
		termbox.Flush()
		ev := termbox.PollEvent()
		// Listen only for keyboard events.
		if ev.Type != termbox.EventKey {
			continue
		}
		switch ev.Key {
		// Name entered.
		case termbox.KeyEnter:
			return name
		// Exit game.
		case ui.CancelKey:
			util.Quit()

		// Erase last character.
		case termbox.KeyBackspace2, termbox.KeyBackspace:
			backspace()
			continue
		// In termbox ev.Ch == 0x00 for ev.Key declared variables. So we need
		// this workaround.
		case termbox.KeySpace:
			addToName(0x20)
			continue
		}
		addToName(ev.Ch)
	}
}
