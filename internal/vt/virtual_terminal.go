package virtual_terminal

import (
	"github.com/charmbracelet/x/vt"
	"github.com/codecrafters-io/shell-tester/internal/screen_state"
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

func (vt *VirtualTerminal) GetScreenState() screen_state.ScreenState {
	cellMatrix := make([][]string, vt.rows)

	for i := 0; i < vt.rows; i++ {
		cellMatrix[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			cellMatrix[i][j] = vt.vt.Cell(j, i).Content
		}
	}

	return screen_state.NewScreenState(cellMatrix, screen_state.CursorPosition{
		RowIndex:    vt.vt.CursorPosition().Y,
		ColumnIndex: vt.vt.CursorPosition().X,
	})
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
