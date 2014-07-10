package text

import (
	"github.com/nsf/termbox-go"
)

type Text struct {
	Text string
	Attr termbox.Attribute
}

func New(str string, attr termbox.Attribute) *Text {
	return &Text{Text: str, Attr: attr}
}
