package virtual_terminal

import (
	"math"
	"strings"

	"github.com/charmbracelet/x/vt"
	"github.com/codecrafters-io/shell-tester/internal/utils"
)

type VirtualTerminal struct {
	vt       *vt.Terminal
	rows     int
	cols     int
	bellChan chan bool
}

func NewStandardVT() *VirtualTerminal {
	// ToDo: This affects performance majorly, improve all functions operating on this
	// Keep a track of when the last row is being written to, panic at that point
	return NewCustomVT(100, 120)
}

func NewCustomVT(rows, cols int) *VirtualTerminal {
	vtInstance := &VirtualTerminal{
		vt:       vt.NewTerminal(cols, rows),
		rows:     rows,
		cols:     cols,
		bellChan: make(chan bool, 1),
	}

	vtInstance.vt.Callbacks.Bell = func() {
		// fmt.Println("🔔 RECEIVED BELL 🔔")
		// Non-blocking send to channel
		select {
		case vtInstance.bellChan <- true:
		default:
		}
	}

	return vtInstance
}

func (vt *VirtualTerminal) BellChannel() chan bool {
	return vt.bellChan
}

func (vt *VirtualTerminal) Close() {
	vt.vt.Close()
}

func (vt *VirtualTerminal) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return vt.vt.Write(p)
}

// GetLastVisibleRowIndex returns the index of the last "visible" row in the terminal
//
// The last visible row is whichever of the following is the last:
// (a) The last non-empty row
// (b) The row with the cursor
func (vt *VirtualTerminal) GetLastVisibleRowIndex() int {
	lastNonEmptyRowIndex := 0
	for i := vt.rows - 1; i >= 0; i-- {
		row := vt.GetRow(i)
		if strings.TrimSpace(strings.Join(row, "")) != "" {
			lastNonEmptyRowIndex = i
			break
		}
	}

	cursorRowIndex, _ := vt.GetCursorPosition()
	return int(math.Max(float64(lastNonEmptyRowIndex), float64(cursorRowIndex)))
}

func (vt *VirtualTerminal) GetScreenState() [][]string {
	if vt == nil {
		return [][]string{}
	}

	cursorRowIndex, cursorColIndex := vt.GetCursorPosition()
	lastVisibleRowIndex := vt.GetLastVisibleRowIndex()

	screenState := make([][]string, lastVisibleRowIndex+1)
	for i := 0; i <= lastVisibleRowIndex; i++ {
		screenState[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			c := vt.vt.Cell(j, i)

			// For the row where the cursor is present, keep all characters until the cursor position
			if i == cursorRowIndex && j < cursorColIndex && c.Content == " " {
				screenState[i][j] = utils.VT_SENTINEL_CHARACTER
			} else {
				screenState[i][j] = c.Content
			}
		}
	}

	return screenState
}

func (vt *VirtualTerminal) GetCursorPosition() (int, int) {
	cursorPosition := vt.vt.CursorPosition()
	cursorRow, cursorCol := cursorPosition.Y, cursorPosition.X
	return cursorRow, cursorCol
}

func (vt *VirtualTerminal) GetMaxColumnCount() int {
	return vt.cols
}

func (vt *VirtualTerminal) GetMaxRowCount() int {
	return vt.rows
}

func (vt *VirtualTerminal) GetRow(row int) []string {
	screenState := make([]string, vt.cols)
	for j := 0; j < vt.cols; j++ {
		c := vt.vt.Cell(j, row)
		screenState[j] = c.Content
	}
	return screenState
}

func BuildCleanedRow(row []string) string {
	result := strings.Join(row, "")
	result = strings.TrimRight(result, " ")

	// VT_SENTINEL_CHARACTER is the representation of " " that we intend to preserve
	result = strings.ReplaceAll(result, utils.VT_SENTINEL_CHARACTER, " ")
	return result
}
