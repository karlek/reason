// Package screen is used to divide the user interface in different parts.
package screen

// Screen is a part of the user interface to be drawn upon.
type Screen struct {
	Width  int
	Height int
	// embedded image.Point
	XOffset int
	YOffset int
}

// AreaScreen is the screen on which areas are drawn.
var AreaScreen Screen
