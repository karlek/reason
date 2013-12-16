// Package action implements actions for creatures.
package action

import (
	"fmt"

	"github.com/karlek/reason/beastiary"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/coord"
	"github.com/karlek/worc/status"
)

///
func Attack(a *area.Area, hero *beastiary.Creature, defender *beastiary.Creature) {
	defender.Hp -= hero.Strength
	status.Print(fmt.Sprintf("You inflict %d damage to %s!", hero.Strength, defender.Name()))
	if defender.Hp <= 0 {
		a.Objects[coord.Coord{defender.X(), defender.Y()}].Pop()
		a.ReDraw(defender.X(), defender.Y())
		status.Print(fmt.Sprintf("You killed %s!", defender.Name()))
	}
}
