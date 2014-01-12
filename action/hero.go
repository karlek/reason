// Package action implements actions for creatures.
package action

import (
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

// HeroTurn listens on user input and then acts on it.
func HeroTurn(sav *save.Save, a *area.Area, hero *creature.Creature) int {
	hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(hero.Hp, hero.MaxHp)
	termbox.Flush()

	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return 0
	}
	switch ev.Ch {
	case ui.LookKey:
		// user wants to look around.
		return look(a, hero)
	case 'm':
		// user wants to try debug function.
		return debug()
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		return pickUp(a, hero)
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		return showInventory(a, hero)
	case ui.DropItemKey:
		// user wants to drop an item.
		return dropItem(a, hero)
	case ui.OpenDoorKey:
		// user wants to open a door.
		return openDoor(a, hero)
	case ui.CloseDoorKey:
		// user wants to close a door.
		return closeDoor(a, hero)
	case ui.QuitKey:
		// user wants to quit game.
		util.Quit()
	case ui.SaveAndQuitKey:
		// user wants to save and exit.
		saveQuit(a, hero, sav)
	}

	// user wants to move hero.
	return heroMovement(ev, a, hero)
}

func debug() int {
	status.Print("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	status.Print("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	return 0
}

func look(a *area.Area, hero *creature.Creature) int {
	Look(*a, hero.X(), hero.Y())
	return 0
}

func saveQuit(a *area.Area, hero *creature.Creature, sav *save.Save) {
	err := sav.Save(*a, *hero)
	if err != nil {
		log.Println(err)
	}
	util.Quit()
}

func closeDoor(a *area.Area, hero *creature.Creature) int {
	actionTaken := CloseDoorNarrative(a, hero.X(), hero.Y())
	if actionTaken {
		return hero.Speed
	}
	return 0
}

func openDoor(a *area.Area, hero *creature.Creature) int {
	actionTaken := OpenDoorNarrative(a, hero.X(), hero.Y())
	if actionTaken {
		return hero.Speed
	}
	return 0
}

func pickUp(a *area.Area, hero *creature.Creature) int {
	actionTaken := PickUpNarrative(a, hero)
	if actionTaken {
		return hero.Speed
	}
	return 0
}

func showInventory(a *area.Area, hero *creature.Creature) int {
	actionTaken := ShowInventory(a, hero)
	if actionTaken {
		return hero.Speed
	}
	return 0
}

func dropItem(a *area.Area, hero *creature.Creature) int {
	actionTaken := DropItem(a, hero)
	if actionTaken {
		return hero.Speed
	}
	return 0
}

func heroMovement(ev termbox.Event, a *area.Area, hero *creature.Creature) int {
	var col *area.Collision
	var err error
	switch ev.Key {
	case ui.MoveUpKey:
		col, err = a.MoveUp(hero)
	case ui.MoveDownKey:
		col, err = a.MoveDown(hero)
	case ui.MoveLeftKey:
		col, err = a.MoveLeft(hero)
	case ui.MoveRightKey:
		col, err = a.MoveRight(hero)
	default:
		return 0
	}
	// Creature moved out of bounds.
	if err != nil {
		return 0
	}
	if col == nil {
		return hero.Speed
	}
	if c, ok := col.S.(*creature.Creature); ok {
		Attack(a, hero, c)
		return hero.Speed
	}
	if fa, ok := col.S.(fauna.Doodad); ok {
		if fa.Name() == "door (closed)" {
			actionTaken := WalkedIntoDoor(a, col.X, col.Y)
			if actionTaken {
				return hero.Speed
			}
		}
	}
	return 0
}
