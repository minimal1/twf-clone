package terminal

import (
	"fmt"
	"unicode/utf8"
)

type EventType int

const (
	KeyPress EventType = iota
)

type Event interface {
	EventType() EventType
}

type KeyPressEvent struct {
	Key  Key
	Rune rune
}

type Key int

const (
	KeyUnknown Key = iota

	KeyEnter
	KeyEsc

	KeyArrowUp
	KeyArrowDown
	KeyArrowRight
	KeyArrowLeft

	KeyTab
	KeyBackspace
	KeyCtrlC
	KeyCtrlD
)

func (e KeyPressEvent) EventType() EventType { return KeyPress }

func parseInputData(data []byte) (Event, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	if len(data) == 1 {
		switch data[0] {
		case 13, 10: // Enter
			return KeyPressEvent{Key: KeyEnter}, nil
		case 27: // ESC
			return KeyPressEvent{Key: KeyEsc}, nil
		case 9: // TAB
			return KeyPressEvent{Key: KeyTab}, nil
		case 8, 127: // Backspace
			return KeyPressEvent{Key: KeyBackspace}, nil
		case 3: // Ctrl+C
			return KeyPressEvent{Key: KeyCtrlC}, nil
		case 4: // Ctrl+D
			return KeyPressEvent{Key: KeyCtrlD}, nil
		}
	}

	if len(data) >= 3 && data[0] == 27 && data[1] == 91 {
		switch data[2] {
		case 65: // UP
			return KeyPressEvent{Key: KeyArrowUp}, nil
		case 66: // DOWN
			return KeyPressEvent{Key: KeyArrowDown}, nil
		case 67: // RIGHT
			return KeyPressEvent{Key: KeyArrowRight}, nil
		case 68: // LEFT
			return KeyPressEvent{Key: KeyArrowLeft}, nil
		default:
			return KeyPressEvent{Key: KeyUnknown}, nil
		}
	}

	r, size := utf8.DecodeRune(data)
	if r == utf8.RuneError && size == 1 {
		return nil, fmt.Errorf("Invalid UTF-8")
	}

	return KeyPressEvent{Rune: r}, nil
}

func (t *Terminal) ReadEvent() (Event, error) {
	buffer := make([]byte, 64)
	n, err := t.in.Read(buffer)
	if err != nil {
		return nil, err
	}

	return parseInputData(buffer[:n])
}
