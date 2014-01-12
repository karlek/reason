// Package action implements actions for creatures.
package action

import (
	"fmt"

	"github.com/karlek/reason/creature"
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

func battleNarrative(a *area.Area, hero *creature.Creature, attacker *creature.Creature) {
	hero.Hp -= attacker.Strength
	status.Print(fmt.Sprintf("You take %d damage from %s!", attacker.Strength, attacker.Name()))
	if hero.Hp <= 0 {
		hero.DrawFOV(a)
		status.Print("You die. Press any key to quit.")
		termbox.Flush()
		termbox.PollEvent()
		util.Quit()
	}
}
