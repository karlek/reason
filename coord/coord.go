// Package coord implements functions for working with points/coordinates on
// planes.
package coord

// Coord is a cartesian coordinate.
type Coord struct {
	X, Y int
}

// Plane is a mathematical plane with fixed height and width.
type Plane struct {
	Width, Height    int
	XOffset, YOffset int
}

// Contains asks if a coordinate c exists in plane p.
func (p Plane) Contains(c Coord) bool {
	return (c.X < p.Width && c.X >= p.XOffset) && (c.Y < p.Height && c.Y >= p.YOffset)
}

/*
	height	y     +---------------------+
			|     |						|
			|     |						|
			|     |						|
			|     |						|
			|     |						|
	yOffset	y	  +---------------------+
			|
------------------x---------------------x--------
			|    xoffset			  width
			|


6+3 = 9
height = 7
*/
