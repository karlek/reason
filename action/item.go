package action

import (
	"fmt"
	"strings"

	"github.com/karlek/reason/area"
	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/ui/status"

	"github.com/nsf/termbox-go"
)

func PickUp(c *creature.Creature, a *area.Area) (actionTaken bool) {
	msg := "There's no item here."
	i, err := c.PickUp(a)
	if i == nil {
		return false
	}

	// Print status message if hero's inventory is full.
	if c.IsHero() {
		if err != nil {
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("%c - %s picked up.", i.Hotkey(), i.String())
		}
	} else {
		msg = fmt.Sprintf("%s picked up %s.", strings.Title(c.Name()), i.String())
	}

	// If the distance to the creature is within the sight radius, print the
	// status message.
	if c.Dist() <= creature.Hero.Sight {
		status.Println(msg, termbox.ColorWhite)
	}

	return true
}
