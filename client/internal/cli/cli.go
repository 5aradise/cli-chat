package cli

import (
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	InputRow    = 29
	MaxInputLen = 118
)

func SaveCursor() {
	fmt.Print("\x1b7")
}

func RestoreCursor() {
	fmt.Print("\x1b8")
}

func MoveTo(row int, column ...int) {
	if len(column) != 0 {
		fmt.Printf("\x1b[%d;%dH", row, column[0])
	} else {
		fmt.Printf("\x1b[%d;H", row)
	}
}

func MoveToInput() {
	MoveTo(InputRow, 2)
}

func ClearConsole() {
	SaveCursor()
	MoveTo(1)
	fmt.Print("\x1b[0J")
	RestoreCursor()
}

func PrintInputFrame() {
	SaveCursor()
	MoveTo(InputRow - 1)
	fmt.Println(strings.Repeat("-", MaxInputLen+2))
	fmt.Println("|" + strings.Repeat(" ", MaxInputLen) + "|")
	fmt.Print(strings.Repeat("-", MaxInputLen+2))
	RestoreCursor()
}

func SafePrintf(printLn *int, format string, a ...any) {
	SaveCursor()
	if *printLn == InputRow-1 {
		ClearConsole()
		PrintInputFrame()
		*printLn = 1
	}
	MoveTo(*printLn)
	*printLn++
	fmt.Printf(format, a...)
	RestoreCursor()
}

func Scan(scanner *bufio.Scanner) (string, int) {
	scanner.Scan()
	input := scanner.Text()
	len := utf8.RuneCountInString(input)
	if len <= MaxInputLen {
		MoveToInput()
		fmt.Print("\x1b[0J")
		fmt.Println(strings.Repeat(" ", MaxInputLen) + "|")
	} else {
		overflow := len / MaxInputLen
		MoveTo(InputRow - overflow-1)
		fmt.Print("\x1b[0J")
		MoveTo(InputRow - 1)
		fmt.Println(strings.Repeat("-", MaxInputLen+2))
		fmt.Println("|" + strings.Repeat(" ", MaxInputLen) + "|")
	}
	fmt.Print(strings.Repeat("-", MaxInputLen+2))
	MoveToInput()
	return input, len
}
