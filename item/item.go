package item

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/model"
)

var Letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ItemModeler interface {
	Itemer
	model.Modeler
}

type Itemer interface {
	IsStackable() bool
	IsEquipable() bool
	GetHotkey() string
	/// TODO(_): rename to GetCount
	GetNum() int
	// TODO(_): implement GetFlavorText
	GetDescription() string
	// TODO(_): remove.
	GetCategory() string
	SetName(string)
	SetDescription(string)
	SetCategory(string)
	// TODO(u): IncCount()
	IncreaseNum(int)
	name.Namer
	// model.Modelable
}

// Item is an object with a name.
type Item struct {
	Itemer
	model.Model
	name        string
	Hotkey      string
	Category    string
	Description string
	Num         int
	Effects     []Effect
}

type Effect struct {
}

type Armor Item

type Jewelery Item

type Weapon Item

type Boots Armor
type Gloves Armor
type Chestwear Armor
type Headgear Armor
type Legwear Armor
type Amulet Jewelery
type Ring Jewelery

func (i Item) IsStackable() bool {
	switch i.Category {
	case "potion":
		return true
	}
	return false
}

func (i Item) IsEquipable() bool {
	switch i.Category {
	case "weapon":
		return true
	}
	return false
}

// Name returns the name of the item.
func (i Item) Name() string {
	return i.name
}

/// GetHotkey
func (i Item) GetHotkey() string {
	return i.Hotkey
}

/// GetNum
func (i Item) GetNum() int {
	return i.Num
}

/// GetDescription
func (i Item) GetDescription() string {
	return i.Description
}

/// GetCategory
func (i Item) GetCategory() string {
	return i.Category
}

/// IncreaseNum
func (i *Item) IncreaseNum(num int) {
	i.Num += num
}

func (i *Item) SetName(n string) {
	i.name = n
}

func (i *Item) SetDescription(desc string) {
	i.Description = desc
}

func (i *Item) SetCategory(cat string) {
	i.Category = cat
}
