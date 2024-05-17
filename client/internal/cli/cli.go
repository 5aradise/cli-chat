package cli

import (
	"bufio"
	"fmt"
)

const InputLn = 29

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
	MoveTo(InputLn, 2)
}

func ClearConsole() {
	SaveCursor()
	MoveTo(1)
	fmt.Print("\x1b[0J")
	RestoreCursor()
}

func PrintInputFrame() {
	SaveCursor()
	MoveTo(InputLn - 1)
	fmt.Print("---------------------------------------------------------------------------------------------------------\n")
	fmt.Print("|                                                                                                       |\n")
	fmt.Print("---------------------------------------------------------------------------------------------------------")
	RestoreCursor()
}

func SafePrintf(printLn *int, format string, a ...any) {
	SaveCursor()
	if *printLn == InputLn-1 {
		ClearConsole()
		PrintInputFrame()
		*printLn = 1
	}
	MoveTo(*printLn)
	*printLn++
	fmt.Printf(format, a...)
	RestoreCursor()
}

func Scan(scanner *bufio.Scanner) string {
	scanner.Scan()
	MoveToInput()
	fmt.Print("\x1b[K")
	fmt.Print("                                                                                                       |")
	MoveToInput()
	return scanner.Text()
}
