// Package action implements actions for creatures.
package action

import (
	"fmt"
	"log"

	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/fauna"
	"github.com/karlek/reason/save"
	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/status"
	"github.com/karlek/reason/util"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/nsf/termbox-go"
)

func Attack(a *area.Area, hero *creature.Creature, defender *creature.Creature) {
	var msg string
	defender.Hp -= hero.Strength
	msg += fmt.Sprintf("You inflict %d damage to %s!", hero.Strength, defender.Name())
	if defender.Hp <= 0 {
		a.Monsters[coord.Coord{defender.X(), defender.Y()}] = nil
		msg += fmt.Sprintf(" You killed %s!", defender.Name())
	}
	status.Print(msg)
}

// HeroTurn listens on user input and then acts on it.
func HeroTurn(sav *save.Save, a *area.Area, hero *creature.Creature) int {
	hero.DrawFOV(a)
	status.Update()
	ui.UpdateHp(hero.Hp, hero.MaxHp)
	termbox.Flush()

	log.Print("")
	// Listen for keystrokes.
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return 0
	}
	switch ev.Ch {
	case ui.LookKey:
		// user wants to look around.
		Look(*a, hero.X(), hero.Y())
		return 0
	case 'm':
		status.Print("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		status.Print("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
		return 0
	case ui.PickUpItemKey:
		// user wants to pick up an item.
		actionTaken := PickUpNarrative(a, hero)
		if actionTaken {
			return hero.Speed
		}
		return 0
	case ui.ShowInventoryKey:
		// user wants to look at inventory.
		actionTaken := ShowInventory(a, hero)
		if actionTaken {
			return hero.Speed
		}
		return 0
	case ui.DropItemKey:
		// user wants to drop an item.
		actionTaken := DropItem(hero, a)
		if actionTaken {
			return hero.Speed
		}
		return 0
	case ui.OpenDoorKey:
		// user wants to open a door.
		actionTaken := OpenDoorNarrative(a, hero.X(), hero.Y())
		if actionTaken {
			return hero.Speed
		}
		return 0
	case ui.CloseDoorKey:
		// user wants to close a door.
		actionTaken := CloseDoorNarrative(a, hero.X(), hero.Y())
		if actionTaken {
			return hero.Speed
		}
		return 0
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
