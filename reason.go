package main

import "fmt"
import "log"
import "math/rand"
import "os"
import "path/filepath"
import "time"

import "github.com/mewkiz/pkg/errorsutil"

import "github.com/karlek/reason/actions"
import "github.com/karlek/reason/beastiary"
import "github.com/karlek/reason/environment"

// ui is initialized with init()
import "github.com/karlek/reason/ui"
import "github.com/karlek/worc/save"
import "github.com/karlek/worc/menu"
import "github.com/karlek/worc/creature"
import "github.com/karlek/worc/terrain"
import "github.com/nsf/termbox-go"

// reason base directory.
var baseDir string

func init() {
	srcDir := "/src/github.com/karlek/reason/"
	for _, goPath := range filepath.SplitList(os.Getenv("GOPATH")) {
		baseDir = goPath + srcDir
		_, err := os.Stat(baseDir)
		if err == nil {
			return
		}
	}
	err := errorsutil.ErrorfColor("unable to locate reason base directory (%q) in GOPATH (%q).", srcDir, os.Getenv("GOPATH"))
	log.Fatalln(err)
}

const (
	// Status messages
	openDoorStr      = "Do you want to open the door? [Y/n]"
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
	err = termbox.Init()
	if err != nil {
		return errorsutil.ErrorfColor("%s", err)
	}
	defer termbox.Close()

	var terrainObjs = []terrain.Terrain{
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Soil,
		environment.Tree,
		// environment.Wall,
		// environment.ClosedDoor,
		// environment.OpenDoor,
	}

	// Tries to load save from path.
	sav, err := save.New(baseDir + "debug.save")
	if err != nil {
		return errorsutil.ErrorfColor("%s", err)
	}

	// Area is where all terrain and object data are stored.
	var a terrain.Area

	// AreaScreen is the active viewport of the area.
	// Used on big areas that are bigger then the screen.
	as := menu.AreaScreen{
		Width:  ui.AreaScreenWidth,
		Height: ui.AreaScreenHeight,
	}

	// The Hero! This is the unit that the user will control
	hero := beastiary.Hero

	// If a save file didn't exist:
	//	generate a new area,
	//	set hero's start location,
	//	spawn some mobs that the hero can debug.
	if !sav.Exists() {
		a = terrain.GenArea(terrainObjs, ui.AreaScreenWidth*2, ui.AreaScreenHeight*2)
		a.Objects[terrain.Coord{X: hero.X, Y: hero.Y}] = hero
		a.SpawnMobs(beastiary.Gobbal)
	} else {
		// If a save file exists, retrive it's data.
		blobs, err := sav.Load()
		if err != nil {
			return errorsutil.ErrorfColor("%s", err)
		}

		// Save files return an empty interface slice.
		// We figure out the type of the stored object
		// and restores it's information.
		for _, object := range blobs {
			switch obj := object.(type) {
			case (*menu.AreaScreen):
				as = *obj
			case (*terrain.Area):
				a = *obj
			default:
				log.Fatalf("Unloaded object: %T\n", obj)
			}
		}

		/// Temporary bug fix (hopefully); when saved, all objects become pointers.
		/// Since all type switches only checks for non-pointers, we convert the types on load.
		// Find the hero object and load his saved progress!
		for coord, object := range a.Objects {
			switch obj := object.(type) {
			case (*creature.Creature):
				a.Objects[coord] = *obj
				if obj.Name == hero.Name {
					hero = a.Objects[coord].(creature.Creature)
				}
			}
		}
	}

	// Draws both terrain and objects to screen.
	a.Draw(as)

	/// Since a.Objects is a map, not the same order
	/// of characters take their turn. Wanted "feature"?
	for {
		for _, object := range a.Objects {
			switch obj := object.(type) {
			case creature.Creature:
				if obj.Name == "Hero" {
					err = HeroAction(&a, &as, &hero, sav)
					if err != nil {
						return errorsutil.ErrorfColor("%s", err)
					}
				} else {
					DoSomething(&obj, &a, &as)
				}
			}
		}
	}
	return nil
}

/// Ugly ugly code
func HeroAction(a *terrain.Area, as *menu.AreaScreen, hero *creature.Creature, sav *save.Save) (err error) {
turn:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			// Open
			case 'o':
				err = actions.Open(a, *hero, *as)
				break turn
			// Look
			case 'l':
				err = actions.Look(a, *hero, *as)
				break turn
			// Quit without saving
			case 'q':
				os.Exit(0)
			}

			switch ev.Key {
			// Save and exit
			case termbox.KeyEsc:
				var wantToSave []interface{}
				wantToSave = append(wantToSave, a)
				wantToSave = append(wantToSave, as)
				err = sav.Save(wantToSave)
				if err != nil {
					return errorsutil.ErrorfColor("%s", err)
				}
				os.Exit(0)

			// Movement
			case termbox.KeyArrowUp:
				err = a.Move(hero, hero.X, hero.Y-1, as)
			case termbox.KeyArrowDown:
				err = a.Move(hero, hero.X, hero.Y+1, as)
			case termbox.KeyArrowLeft:
				err = a.Move(hero, hero.X-1, hero.Y, as)
			case termbox.KeyArrowRight:
				err = a.Move(hero, hero.X+1, hero.Y, as)
			}

			switch err := err.(type) {
			case terrain.MovementError:
				if a.Terrain[err.Y][err.X] == environment.ClosedDoor {
					menu.PrintStatus(openDoorStr)
					actions.WalkedIntoDoor(a, *as, err.X, err.Y)
					break turn
				} else {
					menu.PrintStatus(fmt.Sprintf(pathIsBlockedStr, a.Terrain[err.Y][err.X].Name))
				}
			case nil:
				break turn
			}
		}
	}

	return nil
}

/// Debug function for creature actions
/// If a creature is trapped (see Fig. 1), the loop will never end
/// unless the mob has a passive action.
// Fig. 1
// T#T
// #G+
// ##T
func DoSomething(c *creature.Creature, a *terrain.Area, as *menu.AreaScreen) {
	var err error

mobTurn:
	for {
		switch randInt(0, 5) {
		case 0:
			break mobTurn
		case 1:
			err = a.Move(c, c.X, c.Y-1, as)
		case 2:
			err = a.Move(c, c.X, c.Y+1, as)
		case 3:
			err = a.Move(c, c.X-1, c.Y, as)
		case 4:
			err = a.Move(c, c.X+1, c.Y, as)
		}
		switch err.(type) {
		case nil:
			break mobTurn
		}
	}
}

/// Temp func
func randInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
