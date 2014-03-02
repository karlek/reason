// Reason is a roguelike written on top of worc engine.
package main

import (
	"log"
	"os"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/turn"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"

	"github.com/karlek/worc/area"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

// Main loop function.
func main() {
	err := reason()
	if err != nil {
		log.Fatalln(err)
	}
}

func reason() (err error) {
	// Init graphic library.
	err = termbox.Init()
	if err != nil {
		return errutil.Err(err)
	}

	var sav *save.Save
	var a area.Area
	for {
		Tick(sav, &a)
	}
}

var gameStarted = false
var inited = false
var name string

func mainMenu() {
	const journeyStr = "Do you really need a reason to begin your journey?"
	const nameStr = "What's your name?"

	width, height := termbox.Size()

	journeyX := width/2 - len(journeyStr)/2

	t := ui.NewText(termbox.AttrBold, journeyStr)
	ui.Print(t, journeyX, height/2, 0)
	termbox.Flush()

	ev := termbox.PollEvent()
	if ev.Key == ui.CancelKey {
		termbox.Close()
		os.Exit(0)
	}

	t = ui.NewText(termbox.AttrBold, nameStr)

	offsetX := width/2 - len(nameStr)/2
	offsetY := height/2 + 2

	ui.Print(t, journeyX, offsetY, 0)
	termbox.Flush()

nameLoop:
	for {
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			return
		}
		switch ev.Key {
		// Name entered.
		case termbox.KeyEnter:
			break nameLoop

		// Exit game.
		case ui.CancelKey:
			termbox.Close()
			os.Exit(0)
			return

		// Erase last character.
		case termbox.KeyBackspace2, termbox.KeyBackspace:
			// Remove the last character from name.
			if len(name) > 0 {
				name = name[:len(name)-1]
			}
			ui.ClearLineOffset(offsetY, offsetX)

			t = ui.NewText(termbox.AttrBold+termbox.ColorBlack, name)
			nameT := ui.NewText(termbox.AttrBold, nameStr)
			ui.Print(nameT, journeyX, offsetY, 0)
			ui.Print(t, offsetX+1, offsetY, 0)
			termbox.Flush()
			continue
		}

		// Add new character to name.
		switch ev.Ch {
		default:
			if len(name) > 30 {
				break
			}
			name += string(ev.Ch)
			t = ui.NewText(termbox.AttrBold+termbox.ColorBlack, name)
			ui.Print(t, offsetX+1, offsetY, 0)
			termbox.Flush()
		}
	}
	gameStarted = true
}

func Tick(sav *save.Save, a *area.Area) (err error) {
	if !gameStarted {
		mainMenu()
		return nil
	}

	if !inited {
		// Load or create new game.
		// Load old values or initalize a new area and hero.
		sav, err = initGameSession(a)
		if err != nil {
			return errutil.Err(err)
		}

		// Initalize turn priority queue.
		turn.Init(a)
		inited = true
		status.Print("You will change the world.")
	}
	turn.Proccess(sav, a)
	return nil
}

// initGameLibs initializes creature, fauna and item libraries.
func initGameLibs() (err error) {
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
		newGame(a)
	}
	if err != nil {
		return nil, errutil.Err(err)
	}
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
