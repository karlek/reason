// Package action implements actions for creatures.
package action

import (
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

// HeroTurn listens on user input and then acts on it.
func HeroTurn(sav *save.Save, a *area.Area) int {
	creature.Hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(creature.Hero.Hp, creature.Hero.MaxHp)
	termbox.Flush()

	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return 0
	}
	switch ev.Ch {
	case ui.LookKey:
		// user wants to look around.
		return look(a)
	case 'm':
		// user wants to try debug function.
		return debug()
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		return pickUp(a)
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		return showInventory(a)
	case ui.DropItemKey:
		// user wants to drop an item.
		return dropItem(a)
	case ui.OpenDoorKey:
		// user wants to open a door.
		return openDoor(a)
	case ui.CloseDoorKey:
		// user wants to close a door.
		return closeDoor(a)
	case ui.QuitKey:
		// user wants to quit game.
		util.Quit()
	case ui.SaveAndQuitKey:
		// user wants to save and exit.
		saveQuit(a, sav)
	}

	// user wants to move creature.Hero.
	return HeroMovement(ev, a)
}

func debug() int {
	status.Print("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	status.Print("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	return 0
}

func look(a *area.Area) int {
	Look(*a, creature.Hero.X(), creature.Hero.Y())
	return 0
}

func saveQuit(a *area.Area, sav *save.Save) {
	err := sav.Save(*a)
	if err != nil {
		log.Println(err)
	}
	util.Quit()
}

func closeDoor(a *area.Area) int {
	actionTaken := CloseDoorNarrative(a, creature.Hero.X(), creature.Hero.Y())
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func openDoor(a *area.Area) int {
	actionTaken := OpenDoorNarrative(a, creature.Hero.X(), creature.Hero.Y())
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func pickUp(a *area.Area) int {
	actionTaken := creature.Hero.PickUp(a)
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func showInventory(a *area.Area) int {
	actionTaken := ShowInventory(a)
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func dropItem(a *area.Area) int {
	actionTaken := DropItem(a)
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func HeroMovement(ev termbox.Event, a *area.Area) int {
	var col *area.Collision
	var err error
	switch ev.Key {
	case ui.MoveUpKey:
		col, err = a.MoveUp(&creature.Hero)
	case ui.MoveDownKey:
		col, err = a.MoveDown(&creature.Hero)
	case ui.MoveLeftKey:
		col, err = a.MoveLeft(&creature.Hero)
	case ui.MoveRightKey:
		col, err = a.MoveRight(&creature.Hero)
	default:
		return 0
	}
	// Creature moved out of bounds.
	if err != nil {
		return 0
	}
	if col == nil {
		return creature.Hero.Speed
	}
	if c, ok := col.S.(*creature.Creature); ok {
		creature.Hero.Battle(c, a)
		return creature.Hero.Speed
	}
	if fa, ok := col.S.(terrain.Terrain); ok {
		if fa.Name() == "door (closed)" {
			actionTaken := WalkedIntoDoor(a, col.X, col.Y)
			if actionTaken {
				return creature.Hero.Speed
			}
		}
	}
	return 0
}
