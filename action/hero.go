// Package action implements actions for creatures.
package action

import (
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/state"
	"github.com/karlek/reason/terrain"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/nsf/termbox-go"
)

func MonsterInfo(monsters []*creature.Creature) []ui.MonstInfo {
	mInfo := make([]ui.MonstInfo, len(monsters))
	for i, monst := range monsters {
		mInfo[i] = ui.MonstInfo{Name: monst.Name(), HpLevel: hpLevel(monst), Graphics: monst.Graphic()}
	}
	return mInfo
}

func hpLevel(c *creature.Creature) int {
	switch {
	case float64(c.Hp)/float64(c.MaxHp) > 0.75:
		return 1
	case float64(c.Hp)/float64(c.MaxHp) > 0.5:
		return 2
	case float64(c.Hp)/float64(c.MaxHp) > 0.25:
		return 3
	case float64(c.Hp)/float64(c.MaxHp) >= 0:
		return 4
	}
	return 0
}

// HeroTurn listens on user input and then acts on it.
func HeroTurn(sav *save.Save, a *area.Area) (int, state.State) {
	creature.Hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(creature.Hero.Hp, creature.Hero.MaxHp)
	ui.UpdateMonsterInfo(MonsterInfo(creature.Hero.MonstersInRange(a)))
	termbox.Flush()

	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return 0, state.Wilderness
	}
	switch ev.Ch {
	case ui.LookKey:
		// user wants to look around.
		return look(a), state.Look
	case 'm':
		// user wants to try debug function.
		return debug(a), state.Wilderness
	case '5', 's':
		// user wants to try debug function.
		return creature.Hero.Speed, state.Wilderness
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		return pickUp(a), state.Wilderness
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		return showInventory(a), state.Inventory
	case ui.DropItemKey:
		// user wants to drop an item.
		return dropItem(a), state.Drop
	case ui.OpenDoorKey:
		// user wants to open a door.
		return openDoor(a), state.Open
	case ui.CloseDoorKey:
		// user wants to close a door.
		return closeDoor(a), state.Close
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
