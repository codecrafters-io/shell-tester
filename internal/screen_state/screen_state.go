package screen_state

type CursorPosition struct {
	RowIndex    int
	ColumnIndex int
}

// ScreenState is a representation of the screen state at a given point in time
type ScreenState struct {
	// rawCellMatrix is always of size [rows x cols].
	//
	// Empty cells are represented by " " (space)
	rawCellMatrix [][]string

	cursorPosition CursorPosition
}

func NewScreenState(rawCellMatrix [][]string, cursorPosition CursorPosition) ScreenState {
	return ScreenState{
		rawCellMatrix:  rawCellMatrix,
		cursorPosition: cursorPosition,
	}
}

func (s ScreenState) GetRow(rowIndex int) Row {
	cursorCellIndex := -1

	if s.cursorPosition.RowIndex == rowIndex {
		cursorCellIndex = s.cursorPosition.ColumnIndex
	}

	return Row{
		rawCells:        s.rawCellMatrix[rowIndex],
		cursorCellIndex: cursorCellIndex,
	}
}

func (s ScreenState) GetRowCount() int {
	return len(s.rawCellMatrix)
}

// GetLastLoggableRowIndex returns the index of the last "loggable" row.
//
// This is usually the last non-empty row, but if the cursor is further ahead, it will be the row before the cursor.
func (s ScreenState) GetLastLoggableRowIndex() int {
	lastNonEmptyRowIndex := -1

	for i := 0; i < s.GetRowCount(); i++ {
		if !s.GetRow(i).IsEmpty() {
			lastNonEmptyRowIndex = i
		}
	}

	return max(s.cursorPosition.RowIndex-1, lastNonEmptyRowIndex)
}
