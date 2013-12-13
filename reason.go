// Reason is a roguelike written on top of worc engine.
package main

/// Add y/n dialog for quitting the game.
/// Add saving capabilities.

import (
	"log"

	"github.com/karlek/reason/action"
	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/ui" // loads with init.

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/screen"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

const (
	// Status messages.
	pathIsBlockedStr = "Your path is blocked by %s"
)

// Error wrapper
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
	defer termbox.Close()

	path, err := goutil.SrcDir("github.com/karlek/reason/")
	if err != nil {
		return errutil.Err(err)
	}
	sav, err := save.New(path + "debug.save")
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize beastiary.
	err = beastiary.LoadCreatures()
	if err != nil {
		return errutil.Err(err)
	}

	// Load or create new game.
	var a area.Area
	var hero beastiary.Creature

	// If save exists load old game session.
	if sav.Exists() {
		err = load(sav, &a, &hero)
		if err != nil {
			return errutil.Err(err)
		}
	} else {
		// Otherwise create a new game session.
		newGame(&a, &hero)
	}

	// Draw both terrain and objects to screen.
	a.Draw()

	// Main loop.
	var finished bool
	for !finished {
		finished = playerTurn(sav, &a, &hero)
	}
	return nil
}

// load loads old information from a save file.
func load(sav *save.Save, a *area.Area, hero *beastiary.Creature) error {
	s, err := sav.Load()
	if err != nil {
		return errutil.Err(err)
	}
	*a = s.Area
	*hero = s.Hero

	return nil
}

// newGame initalizes a new game session.
func newGame(a *area.Area, hero *beastiary.Creature) {
	// Create a screen for the area.
	// areaScreen is the active viewport of the area.
	var areaScreen = screen.Screen{
		Width:  ui.AreaScreenWidth,
		Height: ui.AreaScreenHeight,
	}
	*a = gen.Area(areaScreen, areaScreen.Width, areaScreen.Height)

	// Hero starting position.
	*hero = beastiary.Creatures["hero"]
	hero.NewX(2)
	hero.NewY(2)

	a.Objects[coord.Coord{hero.X(), hero.Y()}] = new(area.Stack)
	a.Objects[coord.Coord{hero.X(), hero.Y()}].Push(hero)
	a.Objects[coord.Coord{2, 3}] = new(area.Stack)
	a.Objects[coord.Coord{2, 3}].Push(fauna.Water)
}

// playerTurn listens on user input and then acts on it.
func playerTurn(sav *save.Save, a *area.Area, hero *beastiary.Creature) bool {
	// Listen for keystrokes.
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Ch {
		case 'l':
			// user wants to look around.
			action.Look(hero.X(), hero.Y(), *a)
		case 'q':
			// user wants to quit game.
			return true
		case 'p':
			// user wants to save and exit.
			err := sav.Save(*a, *hero)
			if err != nil {
				log.Println(err)
			}

			return true
		}

		switch ev.Key {
		// Movement.
		case termbox.KeyArrowUp:
			_ = a.MoveUp(hero)
		case termbox.KeyArrowDown:
			_ = a.MoveDown(hero)
		case termbox.KeyArrowLeft:
			_ = a.MoveLeft(hero)
		case termbox.KeyArrowRight:
			_ = a.MoveRight(hero)
		}
	}
	return false
}
