package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/vt"
)

func main() {
	vt := vt.NewTerminal(60, 10)

	_, _ = fmt.Fprintf(vt, "\033[32mHello \033[%dmGolang\033[0m\r\n", 32)

	// echo -e '$ \x1b[6nech\r$ echo \x1b[J'

	// CSI 6n 	DSR 	Device Status Report
	// Reports the cursor position (CPR)
	// by transmitting ESC[n;mR, where n is the row and m is the column.
	_, _ = fmt.Fprintf(vt, "$ \x1b[6n")
	_, _ = fmt.Fprintf(vt, "ech")

	// CSI n J 	ED 	Erase in Display
	// Clears part of the screen.
	// If n is 0 (or missing), clear from cursor to end of screen.
	// If n is 1, clear from cursor to beginning of the screen.
	// If n is 2, clear entire screen (and moves cursor to upper left on DOS ANSI.SYS).
	// If n is 3, clear entire screen and delete all lines saved in the scrollback buffer (this feature was added for xterm and is supported by other terminal applications).
	_, _ = fmt.Fprintf(vt, "\r$ echo \x1b[J")
	logScreenState(vt)
}

func logScreenState(vt *vt.Terminal) {
	screenState := RenderScreenState(vt)
	for _, row := range screenState {
		line := strings.Join(row, "")
		cleanedRow := strings.TrimRight(line, " ")
		if len(cleanedRow) > 0 {
			fmt.Println(line)
		}
	}
}

func RenderScreenState(vt *vt.Terminal) [][]string {
	if vt == nil {
		return [][]string{}
	}

	screenState := make([][]string, vt.Height())
	for i := 0; i < vt.Height(); i++ {
		row := make([]string, vt.Width())
		for j := 0; j < vt.Width(); j++ {
			c := vt.Cell(j, i)
			row[j] = c.Content
		}
		screenState[i] = row
	}

	return screenState
}
