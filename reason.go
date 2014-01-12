/// Issue with saving and saving the current priority queue.
/// Issue with saving and pointers.
// Reason is a roguelike written on top of worc engine.
package main

import (
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/gen"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/turn"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
	"github.com/nsf/termbox-go"
)

// Error wrapper.
func main() {
	err := reason()
	if err != nil {
		log.Fatalln(err)
	}
}

// Main loop function.
func reason() (err error) {
	err = initGameLibs()
	if err != nil {
		return errutil.Err(err)
	}

	var a area.Area
	var hero creature.Creature

	// Load or create new game.
	// Load old values or initalize a new area and hero.
	sav, err := initGameSession(&a, &hero)
	if err != nil {
		return errutil.Err(err)
	}

	// Initalize turn priority queue.
	turn.Init(&a)

	// Main loop.
	for {
		turn.Proccess(sav, &a, &hero)
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

// initGameSession if a load file exists load the old game otherwise
// create a new game session.
func initGameSession(a *area.Area, hero *creature.Creature) (sav *save.Save, err error) {
	path, err := goutil.SrcDir("github.com/karlek/reason/")
	if err != nil {
		return nil, errutil.Err(err)
	}
	sav, err = save.New(path + "debug.save")
	if err != nil {
		return nil, errutil.Err(err)
	}

	// If save exists load old game session.
	// Otherwise create a new game session.
	if sav.Exists() {
		err = load(sav, a, hero)
	} else {
		newGame(a, hero)
	}
	if err != nil {
		return nil, errutil.Err(err)
	}
	return sav, nil
}

// load loads old information from a save file.
func load(sav *save.Save, a *area.Area, hero *creature.Creature) (err error) {
	s, err := sav.Load()
	if err != nil {
		return errutil.Err(err)
	}
	*a = s.Area
	*hero = s.Hero
	return nil
}

// newGame initalizes a new game session.
func newGame(a *area.Area, hero *creature.Creature) error {
	*a = gen.Area(100, 30)
	gen.Mobs(a, 16)
	gen.Items(a, 20)

	// Hero starting position.
	var ok bool
	*hero, ok = creature.Creatures["hero"]
	if !ok {
		return errutil.NewNoPos("unable to locate hero in creatures.")
	}
	hero.SetX(a.Width / 2)
	hero.SetY(a.Height / 2)

	a.Monsters[hero.Coord()] = hero
	return nil
}
