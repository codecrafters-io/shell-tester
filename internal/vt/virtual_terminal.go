package virtual_terminal

import (
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

func (vt *VirtualTerminal) GetScreenState() [][]string {
	if vt == nil {
		return [][]string{}
	}
	cursorPosition := vt.vt.CursorPosition()
	cursorRow, cursorCol := cursorPosition.Y, cursorPosition.X

	// For the row where the cursor is present
	// We intend to keep all characters upto the cursor position
	screenState := make([][]string, vt.rows)
	for i := 0; i < vt.rows; i++ {
		screenState[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			c := vt.vt.Cell(j, i)
			if i == cursorRow && j < cursorCol && c.Content == " " {
				screenState[i][j] = utils.VT_SENTINEL_CHARACTER
			} else {
				screenState[i][j] = c.Content
			}
		}
		// If there is an empty row somewhere in the middle,
		// artificially add a sentinel character
		emptyRowRepresentation := strings.Repeat(" ", vt.cols)
		if i < cursorRow && strings.Join(screenState[i], "") == emptyRowRepresentation {
			screenState[i][0] = utils.VT_SENTINEL_CHARACTER
		}
	}
	return screenState
}

func (vt *VirtualTerminal) GetColumnCount() int {
	return vt.cols
}

func (vt *VirtualTerminal) GetRowCount() int {
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

func (vt *VirtualTerminal) GetRowsTillEnd(startingRow int) [][]string {
	screenState := make([][]string, vt.rows)
	for i := startingRow; i < vt.rows; i++ {
		screenState[i] = vt.GetRow(i)
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
