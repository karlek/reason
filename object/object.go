// Package object adds functionality for stationary object.
package object

import (
	"github.com/karlek/reason/name"

	"github.com/karlek/worc/model"
)

// Object is non-moveable, but drawable object.
type Object struct {
	model.Model
	name.Namer
	name       string
	IsExplored bool
}

// Name returns the name of the object.
func (o Object) Name() string {
	return o.name
}

// New returns a pointer to a copy of o.
func (o Object) New() *Object {
	tmp := o
	return &tmp
}
