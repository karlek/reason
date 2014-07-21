package action

import (
	"github.com/karlek/reason/creature"
	"github.com/karlek/reason/ui"
)

func monsterInfo(monsters []*creature.Creature) []ui.MonstInfo {
	info := make([]ui.MonstInfo, len(monsters))
	for i, monst := range monsters {
		info[i] = ui.MonstInfo{
			Name:     monst.Name(),
			HpLevel:  hpLevel(monst),
			Graphics: monst.Graphic(),
		}
	}
	return info
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
