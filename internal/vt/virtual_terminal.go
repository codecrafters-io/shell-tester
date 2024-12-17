package virtual_terminal

import (
	"strings"

	"github.com/charmbracelet/x/vt"
)

type VirtualTerminal struct {
	vt   *vt.Terminal
	rows int
	cols int
}

func NewStandardVT() *VirtualTerminal {
	// ToDo: This affects performance majorly, improve all functions operating on this
	// Keep a track of when the last row is being written to, panic at that point
	return NewCustomVT(100, 120)
}

func NewCustomVT(rows, cols int) *VirtualTerminal {
	return &VirtualTerminal{
		vt:   vt.NewTerminal(cols, rows),
		rows: rows,
		cols: cols,
	}
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
	screenState := make([][]string, vt.rows)
	for i := 0; i < vt.rows; i++ {
		screenState[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			c := vt.vt.Cell(j, i)
			screenState[i][j] = c.Content
		}
	}
	return screenState
}

func (vt *VirtualTerminal) GetRow(row int, retainColors bool) []string {
	screenState := make([]string, vt.cols)
	for j := 0; j < vt.cols; j++ {
		c := vt.vt.Cell(j, row)
		screenState[j] = string(c.Content)
	}
	return screenState
}

func (vt *VirtualTerminal) GetRowsTillEnd(startingRow int, retainColors bool) [][]string {
	screenState := make([][]string, vt.rows)
	for i := startingRow; i < vt.rows; i++ {
		screenState[i] = vt.GetRow(i, retainColors)
	}
	return screenState
}

func BuildCleanedRow(row []string) string {
	result := strings.Join(row, "")
	result = strings.TrimRight(result, " ")
	return result
}
