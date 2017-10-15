// Package model is something which can be placed on an area and can be drawn.
package model

import (
	"github.com/karlek/reason/coord"

	"github.com/nsf/termbox-go"
)

type Modeler interface {
	X() int
	Y() int
	SetX(int)
	SetY(int)
	Coord() coord.Coord
	SetGraphics(termbox.Cell)
}

// Model is something that is drawed on AreaScreen ontop of an area.
type Model struct {
	Modeler
	Xval     int          // X is x coordinate.
	Yval     int          // Y is y coordinate.
	G        termbox.Cell // G is graphics.
	Pathable bool         // Pathable is a boolean flag if other objects can be placed on top of this model.
}

// Graphic is needed for draw.Drawable.
func (m Model) Graphic() termbox.Cell {
	return m.G
}

// IsPathable is needed for area.Pathable.
func (m Model) IsPathable() bool {
	return m.Pathable
}

// X returns the x value of the current coordinate.
func (m Model) X() int {
	return m.Xval
}

// Y returns the y value of the current coordinate.
func (m Model) Y() int {
	return m.Yval
}

// Coord returns the coordinate.
func (m Model) Coord() coord.Coord {
	return coord.Coord{m.X(), m.Y()}
}

// SetX sets x to the specified coordinate.
func (m *Model) SetX(x int) {
	m.Xval = x
}

// SetY sets y to the specified coordinate.
func (m *Model) SetY(y int) {
	m.Yval = y
}

// SetGraphics sets graphics for the model.
func (m *Model) SetGraphics(g termbox.Cell) {
	m.G = g
}

// SetPathable sets pathability for the model.
func (m *Model) SetPathable(pathable bool) {
	m.Pathable = pathable
}
