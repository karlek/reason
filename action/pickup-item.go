package action

import (
	"strconv"

	"github.com/karlek/reason/creature"

	"github.com/karlek/reason/ui/status"
	"github.com/karlek/worc/area"
)

func PickUpNarrative(a *area.Area, hero *creature.Creature) bool {
	var msg string
	var actionTaken bool
	i := hero.PickUp(a)
	if i == nil {
		msg += "There's no item here."
	} else {
		msg += i.Hotkey + " - "
		if i.IsStackable() {
			if i.Num != 0 {
				msg += strconv.Itoa(i.Num) + " "
			}
		}
		msg += i.Name() + " picked up."
		actionTaken = true
	}
	status.Print(msg)
	return actionTaken
}
