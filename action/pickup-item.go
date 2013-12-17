package action

import (
	"github.com/karlek/reason/beastiary"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/status"
)

///
func PickUpNarrative(a *area.Area, hero *beastiary.Creature) {
	var msg string
	i := hero.PickUp(a)
	if i == nil {
		msg += "There's no item here."
	} else {
		msg += i.Hotkey + " - " + i.Name() + " picked up."
	}
	status.Print(msg)
}
