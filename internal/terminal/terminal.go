package terminal

import (
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	originalState *term.State
	in            *os.File
	out           *os.File
}

func NewTerminal() (*Terminal, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)

	if err != nil {
		return nil, err
	}

	newTerminal := &Terminal{
		in:  tty,
		out: tty,
	}

	return newTerminal, nil
}

func (t *Terminal) EnableRawMode() error {
	originalState, err := term.MakeRaw(int(t.in.Fd()))

	if err != nil {
		return err
	}

	t.originalState = originalState
	return nil
}

func (t *Terminal) DisableRawMode() error {
	if t.originalState == nil {
		return nil
	}

	originalState := t.originalState
	t.originalState = nil
	return term.Restore(int(t.in.Fd()), originalState)
}

func (t *Terminal) GetSize() (width, height int, err error) {
	return term.GetSize(int(t.out.Fd()))
}

func (t *Terminal) Cleanup() error {
	if err := t.DisableRawMode(); err != nil {
		return err
	}

	return t.in.Close()
}
