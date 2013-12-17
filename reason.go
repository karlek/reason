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
	"github.com/karlek/reason/item"
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

	// Initialize items.
	err = item.Load()
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
	ui.Update(hero)

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
	g := *s
	*a = g.Area
	*hero = g.Hero
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
	*a = *gen.Area(areaScreen, areaScreen.Width, areaScreen.Height)
	gen.Mobs(a, 16)
	gen.Items(a, 5)

	// Hero starting position.
	*hero = beastiary.Creatures["hero"]
	hero.NewX(2)
	hero.NewY(2)

	a.Monsters[coord.Coord{hero.X(), hero.Y()}] = hero
}

// nextTurn listens on user input and then acts on it.
func nextTurn(sav *save.Save, a *area.Area, hero *beastiary.Creature) bool {
	// Listen for keystrokes.
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Ch {
		case ui.LookKey:
			// user wants to look around.
			action.Look(*a, hero.X(), hero.Y())
			return false
		case ui.PickUpItemKey:
			// user wants to pick up an item.
			action.PickUpNarrative(a, hero)
			passTime(a, hero)
			return false
		case ui.ShowInventoryKey:
			// user wants to look at inventory.
			actionTaken := action.ShowInventory(a, hero)
			if actionTaken {
				passTime(a, hero)
			}
			ui.AreaScreenRedraw(*a, *hero)
			return false
		case ui.DropItemKey:
			// user wants to drop an item.
			actionTaken := action.DropItem(hero, a)
			if actionTaken {
				passTime(a, hero)
			}
			ui.AreaScreenRedraw(*a, *hero)
			return false
		case ui.OpenDoorKey:
			// user wants to open a door.
			actionTaken := action.OpenDoorNarrative(a, hero.X(), hero.Y())
			if actionTaken {
				passTime(a, hero)
			}
			return false
		case ui.CloseDoorKey:
			// user wants to close a door.
			actionTaken := action.CloseDoorNarrative(a, hero.X(), hero.Y())
			if actionTaken {
				passTime(a, hero)
			}
			return false
		case ui.QuitKey:
			// user wants to quit game.
			return true
		case ui.SaveAndQuitKey:
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
		case ui.MoveUpKey:
			col = a.MoveUp(hero)
		case ui.MoveDownKey:
			col = a.MoveDown(hero)
		case ui.MoveLeftKey:
			col = a.MoveLeft(hero)
		case ui.MoveRightKey:
			col = a.MoveRight(hero)
		}
		if col == nil {
			/// bugs
			passTime(a, hero)
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
	for _, m := range a.Monsters {
		if c, ok := m.(*beastiary.Creature); ok {
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
