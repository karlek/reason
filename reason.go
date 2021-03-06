// Reason is a roguelike written on top of worc engine.
package main

import (
	"fmt"
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/intro"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/object"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/state"
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
	state.Stack.Push(state.Init)

	// Current area.
	a := new(area.Area)
	for {
		err = tick(a)
		if err != nil {
			return err
		}
	}
}

func tick(a *area.Area) (err error) {
	// Save file.
	var sav *save.Save
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
		// Show entry screen, and ask player for character name.
		name := intro.Intro()
		status.Println(fmt.Sprintf("%s. You will change the world.", name), termbox.ColorWhite)
		status.Println("Find Echidna and kill her.", termbox.ColorRed+termbox.AttrBold)
		state.Stack.Push(state.Wilderness)
	case state.Wilderness:
		turn.Proccess(sav, a)
	case state.Inventory:
		fallthrough
	case state.Drop:
		fallthrough
	default:
		state.Stack.Push(state.Wilderness)
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

// initGameLibs initializes creature, fauna and item libraries.
func initGameLibs() (err error) {
	// Init graphic library.
	err = termbox.Init()
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize creatures.
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

	// Initialize Objects.
	err = object.Load()
	if err != nil {
		return errutil.Err(err)
	}
	ui.SetTerminal()

	return nil
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
	// *a = gen.AreaPrim(100, 30)
	*a = gen.Area(100, 30)
	gen.Mobs(a, 16)
	gen.Items(a, 20)
	gen.Objects(a, 50)

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
