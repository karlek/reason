package item

import (
	"fmt"
	"log"
	"strconv"

	"github.com/karlek/reason/item/effect"
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/area"
	"github.com/karlek/worc/model"
)

var Positions string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type DrawItemer interface {
	model.Modeler
	area.DrawPather
	Itemer
}

type Itemer interface {
	SetHotkey(rune)
	Hotkey() rune
	Count() int
	Effects() map[effect.Type]effect.Magnitude
	Rarity() int
	Cat() string
	FlavorText() string
	UseText() string
	SetCount(int)
	fmt.Stringer
	name.Namer
}

// Item rarity ranging from common to artifact.
const (
	Common = iota + 1
	Magical
	Artifact
)

// Item is an object with a name.
type Item struct {
	Itemer
	model.Model
	name     string
	rarity   int
	hotkey   rune
	flavor   string
	use      string
	category string
	count    int
	effects  map[effect.Type]effect.Magnitude
}

// Base types.
type (
	Weapon struct{ Item }
	Potion struct{ Item }
	Tool   struct{ Item }
	Corpse struct{ Item }
	Gold   struct{ Item }
	Scroll struct{ Item }

	Armor     struct{ Item }
	Boots     Armor
	Gloves    Armor
	Chestwear Armor
	Headgear  Armor
	Legwear   Armor

	Jewelery struct{ Item }
	Amulet   Jewelery
	Ring     Jewelery
)

// Name returns the name of the item.
func (i Item) Name() string {
	return i.name
}

// Name returns the name of the item.
func (i Item) Effects() map[effect.Type]effect.Magnitude {
	return i.effects
}

/// Hotkey
func (i Item) Hotkey() rune {
	return i.hotkey
}

/// Num
func (i Item) Count() int {
	return i.count
}

/// FlavorText
func (i Item) FlavorText() string {
	return i.flavor
}

/// UseText
func (i Item) UseText() string {
	return i.use
}

/// Cat
func (i Item) Cat() string {
	return i.category
}

/// Rarity
func (i Item) Rarity() int {
	return i.rarity
}

/// IncreaseNum
func (i *Item) SetCount(n int) {
	i.count = n
}

func (i *Item) SetHotkey(ch rune) {
	i.hotkey = ch
}

func (i *Item) String() string {
	return i.Name()
}

func (i *Potion) String() string {
	return strconv.Itoa(i.Count()) + " " + i.Name()
}

func (i *Scroll) String() string {
	return strconv.Itoa(i.Count()) + " " + i.Name()
}

func IsStackable(i Itemer) bool {
	switch e := i.(type) {
	case *Potion, *Scroll:
		return true
	default:
		log.Printf("%T, not stackable %s", e, e.Name())
		return false
	}
}

func IsEquipable(i Itemer) bool {
	switch i.(type) {
	case *Weapon, *Ring:
		return true
	default:
		return false
	}
}

func IsUsable(i Itemer) bool {
	switch i.(type) {
	case *Potion, *Tool:
		return true
	default:
		return false
	}
}

func IsPermanent(i Itemer) bool {
	switch i.(type) {
	case *Tool:
		return true
	default:
		return false
	}
}
