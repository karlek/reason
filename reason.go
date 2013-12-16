// Reason is a roguelike written on top of worc engine.
package main

/// Add y/n dialog for quitting the game.
/// Add saving capabilities.

import (
	// "fmt"
	"log"
	"math"

	"github.com/karlek/reason/action"
	"github.com/karlek/reason/beastiary"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/ui" // loads with init.

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/screen"
	// "github.com/karlek/worc/status"
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

	// Initialize beastiary.
	err = beastiary.Load()
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize fauna.
	err = fauna.Load()
	if err != nil {
		return errutil.Err(err)
	}

	// Load or create new game.
	var a area.Area
	var hero beastiary.Creature

	// Load or create a new game session.
	sav, err := initGameSession(&a, &hero)
	if err != nil {
		return errutil.Err(err)
	}

	// Initialize and draw the user interface to screen.
	ui.Init(hero)

	// Draw both terrain and objects to screen.
	a.Draw()

	// Main loop.
	var finished bool
	for !finished {
		finished = nextTurn(sav, &a, &hero)
	}
	return nil
}

///
func initGameSession(a *area.Area, hero *beastiary.Creature) (sav *save.Save, err error) {
	path, err := goutil.SrcDir("github.com/karlek/reason/")
	if err != nil {
		return nil, errutil.Err(err)
	}
	sav, err = save.New(path + "debug.save")
	if err != nil {
		return nil, errutil.Err(err)
	}

	// If save exists load old game session.
	if sav.Exists() {
		err = load(sav, a, hero)
		if err != nil {
			return nil, errutil.Err(err)
		}
	} else {
		// Otherwise create a new game session.
		newGame(a, hero)
	}
	return sav, nil
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
	gen.Mobs(a, 14)

	// Hero starting position.
	*hero = beastiary.Creatures["hero"]
	hero.NewX(2)
	hero.NewY(2)

	a.Objects[coord.Coord{hero.X(), hero.Y()}] = new(area.Stack)
	a.Objects[coord.Coord{hero.X(), hero.Y()}].Push(hero)
}

// nextTurn listens on user input and then acts on it.
func nextTurn(sav *save.Save, a *area.Area, hero *beastiary.Creature) bool {
	// Listen for keystrokes.
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Ch {
		case 'l':
			// user wants to look around.
			action.Look(hero.X(), hero.Y(), *a)
			return false
		case 'o':
			actionTaken := action.OpenDoorNarrative(a, hero.X(), hero.Y())
			if actionTaken {
				passTime(a, hero)
			}
			return false
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

		var col *area.Collision
		switch ev.Key {
		// Movement.
		case termbox.KeyArrowUp:
			col = a.MoveUp(hero)
			passTime(a, hero)
		case termbox.KeyArrowDown:
			col = a.MoveDown(hero)
			passTime(a, hero)
		case termbox.KeyArrowLeft:
			col = a.MoveLeft(hero)
			passTime(a, hero)
		case termbox.KeyArrowRight:
			col = a.MoveRight(hero)
			passTime(a, hero)
		}
		if col == nil {
			break
		}
		if c, ok := col.S.(*beastiary.Creature); ok {
			action.Attack(a, hero, c)
			passTime(a, hero)
		}
		if fa, ok := col.S.(fauna.Doodad); ok {
			if fa.Name() == "door (closed)" {
				actionTaken := action.WalkedIntoDoor(a, col.X, col.Y)
				if actionTaken {
					passTime(a, hero)
				}
			}
		}
	}

	ui.UpdateHp(*hero)
	return false
}

///
func passTime(a *area.Area, hero *beastiary.Creature) {
	// Other creatures!
	for _, s := range a.Objects {
		if c, ok := s.Peek().(*beastiary.Creature); ok {
			if c.Name() == "hero" {
				continue
			}

			precise := hero.Speed / c.Speed
			turns := math.Floor(precise)
			reminder := precise - turns
			c.CurSpeed += reminder

			if c.CurSpeed > c.Speed {
				c.CurSpeed -= c.Speed
				turns += 1
			}
			c.Actions(int(turns), a, hero)
		}
	}
}
