package area

import (
	"github.com/karlek/reason/draw"
)

// Pather are objects which can be drawn and answer the question if another
// object can be placed on top of it in the stack.
type Pather interface {
	IsPathable() bool
	SetPathable(bool)
}

type DrawPather interface {
	Pather
	draw.Drawable
}

type DrawIsPather interface {
	IsPathable() bool
	draw.Drawable
}

// Stack is a pile of stuff which can be walked upon or not walked upon.
type Stack []DrawPather

// Peek returns the top most object of the stack without removing it.
func (s *Stack) Peek() DrawPather {
	if s == nil {
		return nil
	}
	if s.Len() == 0 {
		return nil
	}
	return (*s)[s.Len()-1]
}

// Push adds a value ontop of the stack.
func (s *Stack) Push(d DrawPather) {
	*s = append(*s, d)
}

// PopSecond returns the second value from the top.
func (s *Stack) PopSecond() DrawPather {
	if s.Len() < 2 {
		return nil
	}

	tmp := (*s)[s.Len()-2]
	(*s)[s.Len()-2] = (*s)[s.Len()-1]
	*s = (*s)[:s.Len()-1]
	return tmp
}

// Pop returns the latest value.
func (s *Stack) Pop() DrawPather {
	if s.Len() == 0 {
		return nil
	}
	tmp := (*s)[s.Len()-1]
	*s = (*s)[:s.Len()-1]
	return tmp
}

// Len returns the length of the stack.
func (s *Stack) Len() int {
	return len(*s)
}
