// Reason is a roguelike written on top of worc engine.
package main

import (
	// "fmt"
	"log"
	// "math"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/turn"
	"github.com/karlek/reason/ui"
	// "github.com/karlek/reason/ui/status"
	// "github.com/karlek/reason/util"

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
	var hero creature.Creature

	// Load old values or initalize a new area and hero.
	sav, err := loadOrCreateNewGameSession(&a, &hero)
	if err != nil {
		return errutil.Err(err)
	}

	turn.Init(&a)

	// Main loop.
	for {
		turn.Proccess(sav, &a, &hero)
	}
	return nil
}

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

func loadOrCreateNewGameSession(a *area.Area, hero *creature.Creature) (sav *save.Save, err error) {
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
func load(sav *save.Save, a *area.Area, hero *creature.Creature) (err error) {
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
func newGame(a *area.Area, hero *creature.Creature) {
	*a = gen.Area(100, 30)
	gen.Mobs(a, 16)
	gen.Items(a, 20)

	// Hero starting position.
	*hero = creature.Creatures["hero"]
	hero.NewX(ui.Area.Width / 2)
	hero.NewY(ui.Area.Height / 2)

	a.Monsters[coord.Coord{hero.X(), hero.Y()}] = hero
}
