package main

import (
	"fmt"
	"os"

	"github.com/minimal1/twf-clone/internal/filetree"
	"github.com/minimal1/twf-clone/internal/terminal"
)

type FileBrowserApp struct {
	wd string

	terminal    *terminal.Terminal
	fileTree    *filetree.FileTreeImpl
	walker      *filetree.Walker
	currentNode *filetree.TreeNode

	running bool
	startY  int
}

func (app *FileBrowserApp) Initialize() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get wd: %w", err)
	}
	app.wd = wd

	if err := app.readyTerminal(); err != nil {
		return err
	}

	if err := app.readyFileTree(); err != nil {
		return err
	}

	return nil
}

func (app *FileBrowserApp) readyTerminal() error {
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

func (app *FileBrowserApp) readyFileTree() error {
	app.fileTree = filetree.NewFileTree()
	app.walker = filetree.NewWalker(app.fileTree)

	if err := app.fileTree.LoadRoot(app.wd); err != nil {
		return err
	}

	app.currentNode = app.fileTree.GetRoot()

	if err := app.fileTree.ExpandNode(app.currentNode); err != nil {
		return err
	}

	return nil
}

func (app *FileBrowserApp) DrawUI() {
	app.terminal.ClearScreen()
	app.terminal.MoveCursorHome()

	app.drawHeader(app.startY, 2)
	app.drawTree(app.startY+3, 2)
}

func (app *FileBrowserApp) drawHeader(row, col int) {
	app.terminal.WriteColoredAt(row, col, "┌─File Browser────────────────────┐", terminal.ColorCyan)
	app.terminal.WriteColoredAt(row+1, col, fmt.Sprintf("│ Path: %s        │", app.wd), terminal.ColorCyan)
	app.terminal.WriteColoredAt(row+2, col, "└─────────────────────────────────┘", terminal.ColorCyan)
}

func (app *FileBrowserApp) drawTree(row, col int) {
	visibleNodes := app.walker.GetVisibleNodes()

	for i, node := range visibleNodes {
		color := terminal.ColorWhite

		if node == app.currentNode {
			color = terminal.ColorYellow
		}

		leftPadding := 2 * (node.Depth() - 1)

		app.terminal.WriteColoredAt(row+i, col+leftPadding, node.GetDisplayName(), color)
	}
}

func (app *FileBrowserApp) HandleEvent(event terminal.Event) {
	switch e := event.(type) {
	case terminal.KeyPressEvent:
		if e.Key != terminal.KeyUnknown {
			switch e.Key {
			case terminal.KeyCtrlC, terminal.KeyEsc:
				app.running = false
			case terminal.KeyArrowUp:
				app.handleMoveUp()
			case terminal.KeyArrowDown:
				app.handleMoveDown()
			case terminal.KeyEnter:
				app.handleToggleNode()
			}
		} else if e.Rune != 0 {
			switch e.Rune {
			case 'j':
				app.handleMoveDown()
			case 'k':
				app.handleMoveUp()
			case 'l':
				app.handleToggleNode()
			}

		}

	}
}

func (app *FileBrowserApp) handleMoveUp() {
	if node := app.walker.GetPrevVisibleNode(app.currentNode); node != nil {
		app.currentNode = node
	}
}

func (app *FileBrowserApp) handleMoveDown() {
	if node := app.walker.GetNextVisibleNode(app.currentNode); node != nil {
		app.currentNode = node
	}
}

func (app *FileBrowserApp) handleToggleNode() error {
	if !app.currentNode.CanExpand() {
		return nil
	}

	if app.currentNode.Expanded {
		return app.fileTree.CollapseNode(app.currentNode)
	} else {
		return app.fileTree.ExpandNode(app.currentNode)
	}
}

func (app *FileBrowserApp) Cleanup() {
	if app.terminal != nil {
		app.terminal.ShowCursor()
		app.terminal.ExitAltScreen()
		app.terminal.Cleanup()
	}
}

func (app *FileBrowserApp) Run() error {
	if err := app.Initialize(); err != nil {
		return err
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

func main() {
	app := &FileBrowserApp{
		startY:  1,
		running: true,
	}

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
