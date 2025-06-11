package screen_state

import "strings"

type Row struct {
	rawCells        []string
	cursorCellIndex int
}

func (r Row) HasCursor() bool {
	return r.cursorCellIndex != -1
}

func (r Row) String() string {
	if r.HasCursor() {
		// If the cursor is on the line, we need to preserve the spaces before the cursor
		contentsBeforeCursor := strings.Join(r.rawCells[:r.cursorCellIndex], "")
		contentsAfterCursor := strings.TrimRight(strings.Join(r.rawCells[r.cursorCellIndex:], ""), " ")

		return contentsBeforeCursor + contentsAfterCursor
	} else {
		// If the cursor isn't on the line, we can safely trim spaces from the right
		return strings.TrimRight(strings.Join(r.rawCells, ""), " ")
	}
}

func (r Row) IsEmpty() bool {
	return r.String() == ""
}
