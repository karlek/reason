package equipment

import (
	"log"

	"github.com/karlek/reason/item"
	"github.com/karlek/reason/item/effect"
)

type Equipment struct {
	MainHand  *item.Weapon
	OffHand   *item.Weapon
	Head      *item.Headgear
	Amulet    *item.Amulet
	Rings     []*item.Ring
	Boots     *item.Boots
	Gloves    *item.Gloves
	Chestwear *item.Chestwear
	Legwear   *item.Legwear
}

func (eq Equipment) Power() int {
	if eq.MainHand == nil {
		log.Println("no weapon")
		return 0
	}
	if len(eq.MainHand.Effects()) < 1 {
		log.Println("no effect")
		return 0
	}
	if str, ok := eq.MainHand.Effects()[effect.Strength]; ok {
		return int(str)
	}
	return 0
}

func (eq Equipment) Defense() int {
	if eq.MainHand == nil {
		log.Println("no weapon")
		return 0
	}
	if len(eq.MainHand.Effects()) < 1 {
		log.Println("no effect")
		return 0
	}
	if def, ok := eq.MainHand.Effects()[effect.Defense]; ok {
		return int(def)
	}
	return 0
}
