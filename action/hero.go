// Package action implements actions for creatures.
package action

import (
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/inventory"
	"github.com/karlek/reason/item"
	"github.com/karlek/reason/object"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/state"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

// HeroTurn listens on user input and then acts on it.
func HeroTurn(sav *save.Save, a *area.Area) (int, state.State) {
	creature.Hero.DrawFOV(a)
	status.Update()
	ui.Hp(creature.Hero.Hp, creature.Hero.MaxHp)
	ui.Monsters(monsterInfo(creature.Hero.MonstersInRange(a)))
	termbox.Flush()

	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return 0, state.Wilderness
	}
	switch ev.Ch {
	// case 'x':
	// 	sfx.Glitch()
	// 	return 0, state.Wilderness
	case '5', 's':
		// user wants to wait one turn.
		return creature.Hero.Speed, state.Wilderness
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		return pickUp(a), state.Wilderness
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		return showInventory(a), state.Inventory
	// case ui.LookKey:
	// 	// user wants to look around.
	// 	return look(a), state.Look
	// case 'm':
	// 	// user wants to try debug function.
	// 	return debug(a), state.Wilderness
	case ui.DropItemKey:
		// user wants to drop an item.
		return dropItem(a), state.Drop
	// case ui.OpenDoorKey:
	// user wants to open a door.
	// return openDoor(a), state.Open
	// case ui.CloseDoorKey:
	// user wants to close a door.
	// return closeDoor(a), state.Close
	case ui.QuitKey:
		// user wants to quit game.
		util.Quit()
	case ui.SaveAndQuitKey:
		// user wants to save and exit.
		saveQuit(a, sav)
	}

	// user wants to move creature.Hero.
	return HeroMovement(ev, a), state.Wilderness
}

func debug(a *area.Area) int {
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

// Outdated.
// func closeDoor(a *area.Area) int {
// 	actionTaken := CloseDoorNarrative(a, creature.Hero.X(), creature.Hero.Y())
// 	if actionTaken {
// 		return creature.Hero.Speed
// 	}
// 	return 0
// }

// func openDoor(a *area.Area) int {
// 	actionTaken := OpenDoorNarrative(a, creature.Hero.X(), creature.Hero.Y())
// 	if actionTaken {
// 		return creature.Hero.Speed
// 	}
// 	return 0
// }

func pickUp(a *area.Area) int {
	actionTaken := creature.Hero.PickUp(a)
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func showInventory(a *area.Area) int {
	actionTaken := inventory.Show(a)
	if actionTaken {
		return creature.Hero.Speed
	}
	return 0
}

func dropItem(a *area.Area) int {
	actionTaken := inventory.DropItem(a)
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
	stk, ok := a.Items[creature.Hero.Coord()]
	if ok && stk.Len() > 0 {
		if stk.Len() > 1 {
			status.Println("You find a heap of items on the ground:", termbox.ColorWhite)
		} else {
			status.Println("You find a single item on the ground:", termbox.ColorWhite)
		}
		print := "   "
		for _, s := range []area.DrawPather(*stk) {
			i, _ := s.(item.DrawItemer)
			if stk.Len() < 4 {
				print += i.Name() + ", "
			} else {
				print += string(i.Graphic().Ch) + ", "
			}
		}
		status.Println(print[:len(print)-2], termbox.ColorBlack+termbox.AttrBold)
	}
	// Successful movement.
	if col == nil {
		return creature.Hero.Speed
	}
	// Another creature occupied that tile -> battle!
	if c, ok := col.S.(*creature.Creature); ok {
		creature.Hero.Battle(c, a)
		return creature.Hero.Speed
	}
	// Hero walked into an object.
	if obj, ok := a.Objects[col.Coord()].(*object.Object); ok {
		// If the hero walked into a door -> open!
		switch obj.Name() {
		case "door (closed)":
			a.Objects[col.Coord()] = object.Objects["door (open)"].New()
			status.Println("You open the closed door.", termbox.ColorBlack+termbox.AttrBold)
			return creature.Hero.Speed
		}
	}
	return 0
}
