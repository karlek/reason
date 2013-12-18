/// In case whole buffer will be filled with new messages, termbox.PollEvent()?
// Package status implements functions to print status messages.
package status

import (
	"github.com/karlek/reason/ui"

	"github.com/nsf/termbox-go"
)

var messageBuffer = []string{}

// Prints to string to screen taking x coordinate, y coordinate,
// foreground color (attributes) and background color (attributes)
func print(str string, x, y, width int, fg termbox.Attribute, bg termbox.Attribute) {
	// Clears the line from old characters.
	for i := len(str); i < width; i++ {
		termbox.SetCell(x+i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for charOffset, char := range str {
		termbox.SetCell(x+charOffset, y, char, fg, bg)
	}
}

// Print writes a string to the status buffer.
func Print(str string) {
	for {
		if len(str) < ui.Message.Width {
			messageBuffer = append(messageBuffer, str)
			break
		}
		strLen := ui.Message.Width

		messageBuffer = append(messageBuffer, str[:strLen])
		str = str[strLen:]
	}
	Update()
}

func Update() {
	lenBuf := len(messageBuffer)
	lastMessages := messageBuffer
	if lenBuf > ui.Message.Height {
		lastMessages = messageBuffer[lenBuf-ui.Message.Height:]
	}

	for y, msg := range lastMessages {
		print(msg, ui.Message.XOffset, ui.Message.YOffset+y, ui.Message.Width, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
}
