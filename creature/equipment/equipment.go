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

func (eq Equipment) Num() (num int) {
	if eq.MainHand != nil {
		num++
	}
	if eq.OffHand != nil {
		num++
	}
	if eq.Head != nil {
		num++
	}
	if eq.Amulet != nil {
		num++
	}
	if eq.Rings != nil {
		num += len(eq.Rings)
	}
	if eq.Boots != nil {
		num++
	}
	if eq.Gloves != nil {
		num++
	}
	if eq.Chestwear != nil {
		num++
	}
	if eq.Legwear != nil {
		num++
	}
	return num
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
