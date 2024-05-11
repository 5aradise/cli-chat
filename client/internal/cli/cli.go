package cli

import (
	"fmt"
)

type color string

const (
	Reset         = "\x1b[0m"
	Red     color = "\x1b[31m"
	Green   color = "\x1b[32m"
	Yellow  color = "\x1b[33m"
	Blue    color = "\x1b[34m"
	Magenta color = "\x1b[35m"
	Cyan    color = "\x1b[36m"
	Gray    color = "\x1b[37m"
	White   color = "\x1b[97m"
)

func Clear() {
	fmt.Print("\x1b[H\x1b[J")
}

func Color(str string, c color) string {
	return string(c) + str + Reset
}
