package text

import (
	"github.com/nsf/termbox-go"
)

// Text is text displayed in a termbox session.
type Text struct {
	Text string
	Attr termbox.Attribute
	Bg   termbox.Attribute
}

// New text string displayed in termbox with various attributes.
func New(str string, attr termbox.Attribute) *Text {
	return &Text{Text: str, Attr: attr, Bg: termbox.ColorDefault}
}

func NewBg(str string, attr, bg termbox.Attribute) *Text {
	return &Text{Text: str, Attr: attr, Bg: bg}
}
