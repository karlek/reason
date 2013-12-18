package item

import (
	"github.com/karlek/worc/model"
	"github.com/nsf/termbox-go"
)

var Letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Item is an object with a name.
type Item struct {
	M           model.Model
	N           string
	Hotkey      string
	Category    string
	Description string
	FlavorText  string
}

// Name returns the name of the item.
func (i *Item) Name() string {
	return i.N
}

// NewX sets a new x value for the coordinate.
func (i *Item) NewX(x int) {
	i.M.NewX(x)
}

// NewY sets a new y value for the coordinate.
func (i *Item) NewY(y int) {
	i.M.NewY(y)
}

// IsStackable returns whether objects can be stacked ontop of this object.
func (i *Item) IsStackable() bool {
	return i.M.IsStackable()
}

// Graphic returns the graphic data of this object.
func (i *Item) Graphic() termbox.Cell {
	return i.M.Graphic()
}

// X returns the x value of the current coordinate.
func (i *Item) X() int {
	return i.M.X()
}

// Y returns the y value of the current coordinate.
func (i *Item) Y() int {
	return i.M.Y()
}
