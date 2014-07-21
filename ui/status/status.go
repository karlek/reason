/// In case whole buffer will be filled with new messages, termbox.PollEvent()?
// Package status implements functions to print status messages.
package status

import (
	"fmt"
	"math"
	"strings"

	"github.com/karlek/reason/ui"
	"github.com/karlek/reason/ui/text"

	"github.com/nsf/termbox-go"
)

// status is the buffer which all messages are written to.
var status = [][]*text.Text{{}}

// line which current messages are written to; ascending.
var line int

// Print takes a string, an attribute and adds it to the status buffer.
func Print(str string, attr termbox.Attribute) {
	PrintText(text.New(str, attr))
}

// Println takes a string, an attribute and adds it to the status buffer
// suffixed with a newline '\n'.
func Println(str string, attr termbox.Attribute) {
	PrintText(text.New(fmt.Sprintln(str), attr))
}

func PrintTextln(t *text.Text) {
	PrintText(text.New(fmt.Sprintln(t.Text), t.Attr))
}

// Update
func Update() {
	for index := len(status) - 1; index >= 0; index-- {
		if len(status)-index > ui.Message.Height {
			break
		}
		// From 0 to ui.Message.Height, so we don't write all messages on the
		// same row and backwards.
		y := int(math.Abs(float64(index - len(status))))
		printLine(status[index], ui.Message.XOffset, ui.Message.YOffset-y+ui.Message.Height)
	}
}

func printLine(line []*text.Text, x, y int) {
	ui.ClearLine(y)

	// offset foreach character on the line.
	offset := 0
	for _, t := range line {
		for _, char := range t.Text {
			termbox.SetCell(x+offset, y, char, t.Attr, termbox.ColorDefault)
			offset++
		}
	}
}

func lineLen(line []*text.Text) (length int) {
	for _, t := range line {
		length += len(t.Text)
	}
	return
}

func PrintText(t *text.Text) {
	for {
		if line >= len(status) {
			status = append(status, []*text.Text{})
		}
		if lineLen(status[line])+len(t.Text) < ui.Message.Width {
			status[line] = append(status[line], t)

			/// Might be wrong, should remove newline char as well?
			/// What should happen with multiple new lines?
			// If new line.
			if strings.ContainsRune(t.Text, '\n') {
				line++
			}
			break
		}
		strLen := ui.Message.Width - lineLen(status[line])

		status[line] = append(status[line], text.New(t.Text[:strLen], t.Attr))
		line++
		t.Text = t.Text[strLen:]
	}
}
