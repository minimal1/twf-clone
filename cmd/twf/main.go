package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[2J\033[H")

	fmt.Print("터미널 제어 테스트\r\n")
	fmt.Print("아무 키나 누르면 종료됩니다...\r\n")

	var b [1]byte
	os.Stdin.Read(b[:])

	fmt.Print("\033[2J\033[H")
}
