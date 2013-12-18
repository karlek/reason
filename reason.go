// Reason is a roguelike written on top of worc engine.
package main

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
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

// Error wrapper
func main() {
	err := reason()
	if err != nil {
		log.Fatalln(err)
	}
}

func reason() (err error) {
	err = initGameLibs()
	if err != nil {
		return errutil.Err(err)
	}

	// Load or create new game.
	var a area.Area
	var hero beastiary.Creature

	// Load old values or initalize a new area and hero.
	sav, err := loadOrCreateNewGameSession(&a, &hero)
	if err != nil {
		return errutil.Err(err)
	}

	// Main loop.
	for {
		nextTurn(sav, &a, &hero)
	}
	return nil
}

func initGameLibs() (err error) {
	// Init graphic library.
	err = termbox.Init()
	if err != nil {
		return errutil.Err(err)
	}

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

	return nil
}

func loadOrCreateNewGameSession(a *area.Area, hero *beastiary.Creature) (sav *save.Save, err error) {
	// If save exists load old game session.
	path, err := goutil.SrcDir("github.com/karlek/reason/")
	if err != nil {
		return nil, errutil.Err(err)
	}
	sav, err = save.New(path + "debug.save")
	if err != nil {
		return nil, errutil.Err(err)
	}
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
func load(sav *save.Save, a *area.Area, hero *beastiary.Creature) (err error) {
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
	*a = gen.Area(100, 30)
	gen.Mobs(a, 16)
	gen.Items(a, 5)

	// Hero starting position.
	*hero = beastiary.Creatures["hero"]
	hero.NewX(2)
	hero.NewY(2)

	a.Monsters[coord.Coord{hero.X(), hero.Y()}] = hero
}

// nextTurn listens on user input and then acts on it.
func nextTurn(sav *save.Save, a *area.Area, hero *beastiary.Creature) {
	hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(hero.Hp, hero.MaxHp)
	termbox.Flush()

	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return
	}
	switch ev.Ch {
	case ui.LookKey:
		// user wants to look around.
		action.Look(*a, hero.X(), hero.Y())
		return
	case 'm':
		status.Print("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		status.Print("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
		return
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		action.PickUpNarrative(a, hero)
		passTime(a, hero)
		return
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		actionTaken := action.ShowInventory(a, hero)
		if actionTaken {
			passTime(a, hero)
		}
		return
	case ui.DropItemKey:
		// user wants to drop an item.
		actionTaken := action.DropItem(hero, a)
		if actionTaken {
			passTime(a, hero)
		}
		return
	case ui.OpenDoorKey:
		// user wants to open a door.
		actionTaken := action.OpenDoorNarrative(a, hero.X(), hero.Y())
		if actionTaken {
			passTime(a, hero)
		}
		return
	case ui.CloseDoorKey:
		// user wants to close a door.
		actionTaken := action.CloseDoorNarrative(a, hero.X(), hero.Y())
		if actionTaken {
			passTime(a, hero)
		}
		return
	case ui.QuitKey:
		// user wants to quit game.
		util.Quit()
	case ui.SaveAndQuitKey:
		// user wants to save and exit.
		err := sav.Save(*a, *hero)
		if err != nil {
			log.Println(err)
		}
		util.Quit()
	}

	// Movement.
	var col *area.Collision
	switch ev.Key {
	case ui.MoveUpKey:
		col = a.MoveUp(hero)
	case ui.MoveDownKey:
		col = a.MoveDown(hero)
	case ui.MoveLeftKey:
		col = a.MoveLeft(hero)
	case ui.MoveRightKey:
		col = a.MoveRight(hero)
	default:
		return
	}
	if col == nil {
		passTime(a, hero)
		return
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
	hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(hero.Hp, hero.MaxHp)
	termbox.Flush()
}
