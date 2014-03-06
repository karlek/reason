// Reason is a roguelike written on top of worc engine.
package main

import (
	"log"
	"os"
	"unicode"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/state"
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/turn"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

const (
	nameStr    = "What's your name?"
	journeyStr = "Do you really need a reason to begin your journey?"
	logo       = `     ..      ...                               .x+=:.                             ` +
		/* 	  */ "  :~\"8888x :\"%888x                            z`    ^%                            " +
		/* 	  */ ` 8    8888Xf  8888>                              .   <k        u.      u.    u.   ` +
		/* 	  */ `X88x. ?8888k  8888X       .u          u        .@8Ned8"  ...ue888b   x@88k u@88c. ` +
		/* 	  */ `'8888L'8888X  '%88X    ud8888.     us888u.   .@^%8888"   888R Y888r ^"8888""8888" ` +
		/* 	  */ " \"888X 8888X:xnHH(`` :888'8888. .@88 \"8888\" x88:  `)8b.  888R I888>   8888  888R  " +
		/* 	  */ `   ?8~ 8888X X8888   d888 '88%" 9888  9888  8888N=*8888  888R I888>   8888  888R  ` +
		/* 	  */ " -~`   8888> X8888   8888.+\"    9888  9888   %8\"    R88  888R I888>   8888  888R  " +
		/* 	  */ ` :H8x  8888  X8888   8888L      9888  9888    @8Wou 9%  u8888cJ888    8888  888R  ` +
		/* 	  */ " 8888> 888~  X8888   '8888c. .+ 9888  9888  .888888P`    \"*888*P\"    \"*88*\" 8888\" " +
		/* 	  */ " 48\"` '8*~   `8888!`  \"88888%   \"888*\"\"888\" `   ^\"F        'Y\"         \"\"   'Y\"   " +
		/* 	  */ "  ^-==\"\"      `\"\"       \"YP'     ^Y\"   ^Y'                                        "
)

// Main loop function.
func main() {
	err := reason()
	if err != nil {
		log.Fatalln(err)
	}
}

func reason() (err error) {
	var sav *save.Save
	a := new(area.Area)

	state.Stack.Push(state.Init)
	for {
		tick(sav, a)
	}
}

var name string

func mainMenu() {
	printJourney()
	printLogo()
	printEnterSign()
	waitForInput()
	removeEnterSign()
	printNameQuestion()
	askName()
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

func askName() {
	for {
		termbox.Flush()
		ev := termbox.PollEvent()
		// Listen only for keyboard events.
		if ev.Type != termbox.EventKey {
			return
		}
		switch ev.Key {
		// Name entered.
		case termbox.KeyEnter:
			return
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

func tick(sav *save.Save, a *area.Area) (err error) {
	switch state.Stack.Pop() {
	case state.Init:
		// Load or create new game.
		// Load old values or initalize a new area and hero.
		sav, err = initGameSession(a)
		if err != nil {
			return errutil.Err(err)
		}
		state.Stack.Push(state.Intro)
	case state.Intro:
		mainMenu()
		status.Print("You will change the world.")
		state.Stack.Push(state.Wilderness)
	case state.Wilderness:
		turn.Proccess(sav, a)
	// case state.Inventory:
	default:
		state.Stack.Push(state.Wilderness)
	}
	return nil
}

// initGameLibs initializes creature, fauna and item libraries.
func initGameLibs() (err error) {
	// Init graphic library.
	err = termbox.Init()
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize creature.
	err = creature.Load()
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize fauna.
	err = terrain.Load()
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize items.
	err = item.Load()
	if err != nil {
		return errutil.Err(err)
	}

	return nil
}

// initGameSession if a load file exists load the old game otherwise
// create a new game session.
func initGameSession(a *area.Area) (sav *save.Save, err error) {
	path, err := goutil.SrcDir("github.com/karlek/reason/")
	if err != nil {
		return nil, errutil.Err(err)
	}
	sav, err = save.New(path + "debug.save")
	if err != nil {
		return nil, errutil.Err(err)
	}

	err = initGameLibs()
	if err != nil {
		return nil, errutil.Err(err)
	}

	// If save exists load old game session.
	// Otherwise create a new game session.
	if sav.Exists() {
		err = load(sav, a)
	} else {
		err = newGame(a)
	}
	if err != nil {
		return nil, errutil.Err(err)
	}
	// Initalize turn priority queue.
	turn.Init(a)

	return sav, nil
}

// load loads old information from a save file.
func load(sav *save.Save, a *area.Area) (err error) {
	s, err := sav.Load()
	if err != nil {
		return errutil.Err(err)
	}
	*a = s.Area
	creature.Hero = s.Hero
	return nil
}

// newGame initalizes a new game session.
func newGame(a *area.Area) error {
	*a = gen.Area(100, 30)
	gen.Mobs(a, 16)
	gen.Items(a, 20)

	// Hero starting position.
	var ok bool
	creature.Hero, ok = creature.Beastiary["hero"]
	if !ok {
		return errutil.NewNoPos("unable to locate hero in creatures.")
	}
	creature.Hero.SetX(a.Width / 2)
	creature.Hero.SetY(a.Height / 2)

	a.Monsters[creature.Hero.Coord()] = &creature.Hero
	return nil
}
