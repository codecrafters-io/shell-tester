package shell_executable

import (
	"github.com/Edgaru089/vterm"
	"github.com/gookit/color"
)

const VT_SENTINEL_CHARACTER = "."

// TODO: See if we can remove this wrapper entirely, and only add helpers to fetch rows/data from the vterm instance directly
type VirtualTerminal struct {
	vt   *vterm.VTerm
	rows int
	cols int
}

func NewStandardVT() *VirtualTerminal {
	return NewCustomVT(12, 80)
}

func NewCustomVT(rows, cols int) *VirtualTerminal {
	return &VirtualTerminal{
		vt:   vterm.New(rows, cols),
		rows: rows,
		cols: cols,
	}
}

func (vt *VirtualTerminal) Close() {
	vt.vt.Free()
}

func (vt *VirtualTerminal) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return vt.vt.Write(p)
}

func (vt *VirtualTerminal) GetScreenState(retainColors bool) [][]string {
	screenState := make([][]string, vt.rows)
	for i := 0; i < vt.rows; i++ {
		screenState[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			c := vt.vt.CellAt(i, j)
			if len(c.Chars) > 0 {
				screenState[i][j] = string(c.Chars)
			} else {
				screenState[i][j] = VT_SENTINEL_CHARACTER
			}
		}
	}
	return screenState
}

func (vt *VirtualTerminal) GetRow(row int, retainColors bool) []string {
	screenState := make([]string, vt.cols)
	for j := 0; j < vt.cols; j++ {
		c := vt.vt.CellAt(row, j)
		fr, fg, fb := vt.vt.ConvertRGB(&c.Foreground)
		br, bg, bb := vt.vt.ConvertRGB(&c.Background)
		style := getForegroundBackgroundStyleFromRGB(fr, fg, fb, br, bg, bb)
		if len(c.Chars) > 0 {
			if retainColors {
				screenState[j] = style.Sprintf("%c", c.Chars[0])
			} else {
				screenState[j] = string(c.Chars)
			}
		} else {
			screenState[j] = VT_SENTINEL_CHARACTER
		}
	}
	return screenState
}

func (vt *VirtualTerminal) GetRowsTillEnd(startingRow int, retainColors bool) [][]string {
	screenState := make([][]string, vt.rows)
	for i := startingRow; i < vt.rows; i++ {
		screenState[i] = make([]string, vt.cols)
		for j := 0; j < vt.cols; j++ {
			c := vt.vt.CellAt(i, j)
			fr, fg, fb := vt.vt.ConvertRGB(&c.Foreground)
			br, bg, bb := vt.vt.ConvertRGB(&c.Background)
			style := getForegroundBackgroundStyleFromRGB(fr, fg, fb, br, bg, bb)
			if len(c.Chars) > 0 {
				if retainColors {
					screenState[i][j] = style.Sprintf("%c", c.Chars[0])
				} else {
					screenState[i][j] = string(c.Chars)
				}
			} else {
				screenState[i][j] = VT_SENTINEL_CHARACTER
			}
		}
	}
	return screenState
}

func getForegroundBackgroundStyleFromRGB(fr, fg, fb, br, bg, bb uint8) *color.RGBStyle {
	style := color.NewRGBStyle(
		color.RGB(fr, fg, fb), // Foreground color
		color.RGB(br, bg, bb), // Background color
	)
	return style
}