package main

import (
	"fmt"

	"github.com/minimal1/twf-clone/internal/terminal"
)

type DemoApp struct {
	terminal *terminal.Terminal
	cursorX  int
	cursorY  int
	lastKey  string
	running  bool
}

func main() {
	app := &DemoApp{
		cursorX: 2,
		cursorY: 15,
		lastKey: "None",
		running: true,
	}

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (app *DemoApp) Run() error {
	if err := app.Initialize(); err != nil {
		return nil
	}

	defer app.Cleanup()

	app.DrawUI()

	for app.running {
		event, err := app.terminal.ReadEvent()

		if err != nil {
			return err
		}

		app.HandleEvent(event)
		app.DrawUI()
	}

	return nil
}

func (app *DemoApp) Initialize() error {
	terminal, err := terminal.NewTerminal()
	if err != nil {
		return fmt.Errorf("failed to create terminal: %w", err)
	}
	app.terminal = terminal

	if err := app.terminal.EnableRawMode(); err != nil {
		return fmt.Errorf("failed to enable raw mode: %w", err)
	}

	if err := app.terminal.EnterAltScreen(); err != nil {
		return fmt.Errorf("failed to enter alt screen: %w", err)
	}

	if err := app.terminal.HideCursor(); err != nil {
		return fmt.Errorf("failed to hide cursor: %w", err)
	}

	if err := app.terminal.ClearScreen(); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}

	return nil
}

func (app *DemoApp) DrawUI() {
	app.terminal.ClearScreen()
	app.terminal.MoveCursorHome()

	app.terminal.WriteColoredAt(1, 2, "┌─────────────────────────────────┐", terminal.ColorCyan)
	app.terminal.WriteColoredAt(2, 2, "│        TUI Terminal Demo        │", terminal.ColorCyan)
	app.terminal.WriteColoredAt(3, 2, "└─────────────────────────────────┘", terminal.ColorCyan)

	app.terminal.WriteColoredAt(5, 2, "Controls:", terminal.ColorWhite)
	app.terminal.WriteColoredAt(6, 2, "- Arrow Keys: Move Cursor", terminal.ColorWhite)
	app.terminal.WriteColoredAt(7, 2, "- ESC or Ctrl+C: Exit", terminal.ColorWhite)
	app.terminal.WriteColoredAt(8, 2, "- Any key: Display key info", terminal.ColorWhite)

	app.terminal.WriteColoredAt(10, 2, "┌─Status──────────────────────────┐", terminal.ColorGreen)
	app.terminal.WriteColoredAt(11, 2, fmt.Sprintf("│Cursor: (%d, %d)                 │", app.cursorX, app.cursorY), terminal.ColorGreen)
	app.terminal.WriteColoredAt(12, 2, fmt.Sprintf("│Last Key: %s                     │", app.lastKey), terminal.ColorGreen)
	app.terminal.WriteColoredAt(13, 2, "└─────────────────────────────────┘", terminal.ColorGreen)

	app.terminal.WriteColoredAt(app.cursorY, app.cursorX, "*", terminal.ColorRed)
}

func (app *DemoApp) HandleEvent(event terminal.Event) {
	switch e := event.(type) {
	case terminal.KeyPressEvent:
		if e.Key != terminal.KeyUnknown {
			switch e.Key {
			case terminal.KeyCtrlC, terminal.KeyEsc:
				app.running = false
				app.lastKey = "Exit"
			case terminal.KeyArrowUp:
				if app.cursorY > 1 {
					app.cursorY--
				}
				app.lastKey = "Up"
			case terminal.KeyArrowDown:
				if app.cursorY < 30 {
					app.cursorY++
				}
				app.lastKey = "Down"
			case terminal.KeyArrowRight:
				if app.cursorX < 30 {
					app.cursorX++
				}
				app.lastKey = "Right"
			case terminal.KeyArrowLeft:
				if app.cursorX > 1 {
					app.cursorX--
				}
				app.lastKey = "Left"
			case terminal.KeyTab:
				app.lastKey = "Tab"
			case terminal.KeyBackspace:
				app.lastKey = "Backspace"
			}
		} else if e.Rune != 0 {
			app.lastKey = fmt.Sprintf("'%c'", e.Rune)
		}
	}
}

func (app *DemoApp) Cleanup() {
	if app.terminal != nil {
		app.terminal.ShowCursor()
		app.terminal.ExitAltScreen()
		app.terminal.Cleanup()
	}
}
