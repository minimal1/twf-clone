package terminal

import "fmt"

const (
	// 화면 제어
	ClearScreen     = "\x1b[2J"
	ClearLine       = "\x1b[2K"
	ClearFromCursor = "\x1b[0J"

	// 커서 제어
	CursorHome = "\x1b[H"
	CursorHide = "\x1b[?25l"
	CursorShow = "\x1b[?25h"

	// 스크린 모드
	AltScreenOn  = "\x1b[?1049h"
	AltScreenOff = "\x1b[?1049l"
)

// 화면 제어
func (t *Terminal) ClearScreen() error {
	_, err := t.out.Write([]byte(ClearScreen))
	return err
}
func (t *Terminal) ClearLine() error {
	_, err := t.out.Write([]byte(ClearLine))
	return err
}
func (t *Terminal) ClearFromCursor() error {
	_, err := t.out.Write([]byte(ClearFromCursor))
	return err
}

// 스크린 제어
func (t *Terminal) EnterAltScreen() error {
	_, err := t.out.Write([]byte(AltScreenOn))
	return err
}
func (t *Terminal) ExitAltScreen() error {
	_, err := t.out.Write([]byte(AltScreenOff))
	return err
}

// Cursor 제어
func (t *Terminal) MoveCursorHome() error {
	_, err := t.out.Write([]byte(CursorHome))
	return err
}
func (t *Terminal) MoveCursorTo(row, col int) error {
	// ANSI: \x1b[{row};{col}H
	sequence := fmt.Sprintf("\x1b[%d;%dH", row, col)
	_, err := t.out.Write([]byte(sequence))
	return err
}
func (t *Terminal) HideCursor() error {
	_, err := t.out.Write([]byte(CursorHide))
	return err
}
func (t *Terminal) ShowCursor() error {
	_, err := t.out.Write([]byte(CursorShow))
	return err
}

type Color string

const (
	ColorReset  Color = "\x1b[0m"
	ColorBlack  Color = "\x1b[30m"
	ColorRed    Color = "\x1b[31m"
	ColorGreen  Color = "\x1b[32m"
	ColorYellow Color = "\x1b[33m"
	ColorBlue   Color = "\x1b[34m"
	ColorPurple Color = "\x1b[35m"
	ColorCyan   Color = "\x1b[36m"
	ColorWhite  Color = "\x1b[37m"
)

type Style string

const (
	StyleBold      Style = "\x1b[1m"
	StyleDim       Style = "\x1b[2m"
	StyleUnderline Style = "\x1b[4m"
)

func (t *Terminal) WriteColored(text string, color Color) error {
	coloredText := string(color) + text + string(ColorReset)
	_, err := t.out.Write([]byte(coloredText))
	return err
}

func (t *Terminal) WriteColoredAt(row, col int, text string, color Color) error {
	if err := t.MoveCursorTo(row, col); err != nil {
		return err
	}

	return t.WriteColored(text, color)
}
