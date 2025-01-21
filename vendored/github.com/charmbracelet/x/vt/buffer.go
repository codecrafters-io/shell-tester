package vt

import (
	"github.com/charmbracelet/x/wcwidth"
	"github.com/rivo/uniseg"
)

// NewCell returns a new cell. This is a convenience function that initializes a
// new cell with the given content. The cell's width is determined by the
// content using [wcwidth.RuneWidth].
func NewCell(r rune) *Cell {
	return &Cell{Content: string(r), Width: wcwidth.RuneWidth(r)}
}

// NewGraphemeCell returns a new cell. This is a convenience function that
// initializes a new cell with the given content. The cell's width is determined
// by the content using [uniseg.FirstGraphemeClusterInString].
// This is used when the content is a grapheme cluster i.e. a sequence of runes
// that form a single visual unit.
// This will only return the first grapheme cluster in the string. If the
// string is empty, it will return an empty cell with a width of 0.
func NewGraphemeCell(s string) *Cell {
	c, _, w, _ := uniseg.FirstGraphemeClusterInString(s, -1)
	return &Cell{Content: c, Width: w}
}

var blankCell = Cell{
	Content: " ",
	Width:   1,
}

// Line represents a line in the terminal.
// A nil cell represents an blank cell, a cell with a space character and a
// width of 1.
// If a cell has no content and a width of 0, it is a placeholder for a wide
// cell.
type Line []*Cell

// Width returns the width of the line.
func (l Line) Width() int {
	return len(l)
}

// Buffer is a 2D grid of cells representing a screen or terminal.
type Buffer struct {
	lines []Line
}

// NewBuffer creates a new buffer with the given width and height.
// This is a convenience function that initializes a new buffer and resizes it.
func NewBuffer(width int, height int) *Buffer {
	b := new(Buffer)
	b.Resize(width, height)
	return b
}

// Cell implements Screen.
func (b *Buffer) Cell(x int, y int) *Cell {
	if y < 0 || y >= len(b.lines) {
		return nil
	}
	if x < 0 || x >= b.lines[y].Width() {
		return nil
	}

	c := b.lines[y][x]
	if c == nil {
		newCell := blankCell
		return &newCell
	}

	return c
}

// Draw implements Screen.
func (b *Buffer) Draw(x int, y int, c Cell) bool {
	return b.SetCell(x, y, &c)
}

// maxCellWidth is the maximum width a terminal cell can get.
const maxCellWidth = 4

// SetCell sets the cell at the given x, y position.
func (b *Buffer) SetCell(x, y int, c *Cell) bool {
	return b.setCell(x, y, c, true)
}

// setCell sets the cell at the given x, y position. This will always clone and
// allocates a new cell if c is not nil.
func (b *Buffer) setCell(x, y int, c *Cell, clone bool) bool {
	if y < 0 || y >= len(b.lines) {
		return false
	}
	width := b.lines[y].Width()
	if x < 0 || x >= width {
		return false
	}

	// When a wide cell is partially overwritten, we need
	// to fill the rest of the cell with space cells to
	// avoid rendering issues.
	prev := b.Cell(x, y)
	if prev != nil && prev.Width > 1 {
		// Writing to the first wide cell
		for j := 0; j < prev.Width && x+j < b.lines[y].Width(); j++ {
			newCell := *prev
			newCell.Content = " "
			newCell.Width = 1
			b.lines[y][x+j] = &newCell
		}
	} else if prev != nil && prev.Width == 0 {
		// Writing to wide cell placeholders
		for j := 1; j < maxCellWidth && x-j >= 0; j++ {
			wide := b.Cell(x-j, y)
			if wide != nil && wide.Width > 1 {
				for k := 0; k < wide.Width; k++ {
					newCell := *wide
					newCell.Content = " "
					newCell.Width = 1
					b.lines[y][x-j+k] = &newCell
				}
				break
			}
		}
	}

	if clone && c != nil {
		// Clone the cell if not nil.
		newCell := *c
		c = &newCell
	}

	if c != nil && x+c.Width > width {
		// If the cell is too wide, we write blanks with the same style.
		for i := 0; i < c.Width && x+i < width; i++ {
			newCell := *c
			newCell.Content = " "
			newCell.Width = 1
			b.lines[y][x+i] = &newCell
		}
	} else {
		b.lines[y][x] = c

		// Mark wide cells with an empty cell zero width
		// We set the wide cell down below
		if c != nil && c.Width > 1 {
			for j := 1; j < c.Width && x+j < b.lines[y].Width(); j++ {
				var wide Cell
				b.lines[y][x+j] = &wide
			}
		}
	}

	return true
}

// Height implements Screen.
func (b *Buffer) Height() int {
	return len(b.lines)
}

// Width implements Screen.
func (b *Buffer) Width() int {
	if len(b.lines) == 0 {
		return 0
	}
	return b.lines[0].Width()
}

// Bounds returns the bounds of the buffer.
func (b *Buffer) Bounds() Rectangle {
	return Rect(0, 0, b.Width(), b.Height())
}

// Resize resizes the buffer to the given width and height.
func (b *Buffer) Resize(width int, height int) {
	if width == 0 || height == 0 {
		b.lines = nil
		return
	}

	if width > b.Width() {
		line := make(Line, width-b.Width())
		for i := range b.lines {
			b.lines[i] = append(b.lines[i], line...)
		}
	} else if width < b.Width() {
		for i := range b.lines {
			b.lines[i] = b.lines[i][:width]
		}
	}

	if height > len(b.lines) {
		for i := len(b.lines); i < height; i++ {
			b.lines = append(b.lines, make(Line, width))
		}
	} else if height < len(b.lines) {
		b.lines = b.lines[:height]
	}
}

// fill fills the buffer with the given cell and rectangle.
func (b *Buffer) fill(c *Cell, rect Rectangle) {
	cellWidth := 1
	if c != nil {
		cellWidth = c.Width
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x += cellWidth {
			b.setCell(x, y, c, true) //nolint:errcheck
		}
	}
}

// Fill fills the buffer with the given cell and rectangle.
func (b *Buffer) Fill(c *Cell, rects ...Rectangle) {
	if len(rects) == 0 {
		b.fill(c, b.Bounds())
		return
	}
	for _, rect := range rects {
		b.fill(c, rect)
	}
}

// Clear clears the buffer with space cells and rectangle.
func (b *Buffer) Clear(rects ...Rectangle) {
	if len(rects) == 0 {
		b.fill(nil, b.Bounds())
		return
	}
	for _, rect := range rects {
		b.fill(nil, rect)
	}
}

// InsertLine inserts n lines at the given line position, with the given
// optional cell, within the specified rectangles. If no rectangles are
// specified, it inserts lines in the entire buffer. Only cells within the
// rectangle's horizontal bounds are affected. Lines are pushed out of the
// rectangle bounds and lost. This follows terminal [ansi.IL] behavior.
func (b *Buffer) InsertLine(y, n int, c *Cell, rects ...Rectangle) {
	if len(rects) == 0 {
		b.insertLineInRect(y, n, c, b.Bounds())
	}
	for _, rect := range rects {
		b.insertLineInRect(y, n, c, rect)
	}
}

// insertLineInRect inserts new lines at the given line position, with the
// given optional cell, within the rectangle bounds. Only cells within the
// rectangle's horizontal bounds are affected. Lines are pushed out of the
// rectangle bounds and lost. This follows terminal [ansi.IL] behavior.
func (b *Buffer) insertLineInRect(y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// Limit number of lines to insert to available space
	if y+n > rect.Max.Y {
		n = rect.Max.Y - y
	}

	// Move existing lines down within the bounds
	for i := rect.Max.Y - 1; i >= y+n; i-- {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// We don't need to clone c here because we're just moving lines down.
			// b.lines[i][x] = b.lines[i-n][x]
			b.setCell(x, i, b.lines[i-n][x], false)
		}
	}

	// Clear the newly inserted lines within bounds
	for i := y; i < y+n; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.setCell(x, i, c, true)
		}
	}
}

// deleteLineInRect deletes lines at the given line position, with the given
// optional cell, within the rectangle bounds. Only cells within the
// rectangle's bounds are affected. Lines are shifted up within the bounds and
// new blank lines are created at the bottom. This follows terminal [ansi.DL]
// behavior.
func (b *Buffer) deleteLineInRect(y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// Limit deletion count to available space in scroll region
	if n > rect.Max.Y-y {
		n = rect.Max.Y - y
	}

	// Shift cells up within the bounds
	for dst := y; dst < rect.Max.Y-n; dst++ {
		src := dst + n
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// We don't need to clone c here because we're just moving cells up.
			// b.lines[dst][x] = b.lines[src][x]
			b.setCell(x, dst, b.lines[src][x], false)
		}
	}

	// Fill the bottom n lines with blank cells
	for i := rect.Max.Y - n; i < rect.Max.Y; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.setCell(x, i, c, true)
		}
	}
}

// DeleteLine deletes n lines at the given line position, with the given
// optional cell, within the specified rectangles. If no rectangles are
// specified, it deletes lines in the entire buffer.
func (b *Buffer) DeleteLine(y, n int, c *Cell, rects ...Rectangle) {
	if len(rects) == 0 {
		b.deleteLineInRect(y, n, c, b.Bounds())
		return
	}
	for _, rect := range rects {
		b.deleteLineInRect(y, n, c, rect)
	}
}

// InsertCell inserts new cells at the given position, with the given optional
// cell, within the specified rectangles. If no rectangles are specified, it
// inserts cells in the entire buffer. This follows terminal [ansi.ICH]
// behavior.
func (b *Buffer) InsertCell(x, y, n int, c *Cell, rects ...Rectangle) {
	if len(rects) == 0 {
		b.insertCellInRect(x, y, n, c, b.Bounds())
		return
	}
	for _, rect := range rects {
		b.insertCellInRect(x, y, n, c, rect)
	}
}

// insertCellInRect inserts new cells at the given position, with the given
// optional cell, within the rectangle bounds. Only cells within the
// rectangle's bounds are affected, following terminal [ansi.ICH] behavior.
func (b *Buffer) insertCellInRect(x, y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// Limit number of cells to insert to available space
	if x+n > rect.Max.X {
		n = rect.Max.X - x
	}

	// Move existing cells within rectangle bounds to the right
	for i := rect.Max.X - 1; i >= x+n && i-n >= rect.Min.X; i-- {
		// We don't need to clone c here because we're just moving cells to the
		// right.
		// b.lines[y][i] = b.lines[y][i-n]
		b.setCell(i, y, b.lines[y][i-n], false)
	}

	// Clear the newly inserted cells within rectangle bounds
	for i := x; i < x+n && i < rect.Max.X; i++ {
		b.setCell(i, y, c, true)
	}
}

// DeleteCell deletes cells at the given position, with the given optional
// cell, within the specified rectangles. If no rectangles are specified, it
// deletes cells in the entire buffer. This follows terminal [ansi.DCH]
// behavior.
func (b *Buffer) DeleteCell(x, y, n int, c *Cell, rects ...Rectangle) {
	if len(rects) == 0 {
		b.deleteCellInRect(x, y, n, c, b.Bounds())
		return
	}
	for _, rect := range rects {
		b.deleteCellInRect(x, y, n, c, rect)
	}
}

// deleteCellInRect deletes cells at the given position, with the given
// optional cell, within the rectangle bounds. Only cells within the
// rectangle's bounds are affected, following terminal [ansi.DCH] behavior.
func (b *Buffer) deleteCellInRect(x, y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// Calculate how many positions we can actually delete
	remainingCells := rect.Max.X - x
	if n > remainingCells {
		n = remainingCells
	}

	// Shift the remaining cells to the left
	for i := x; i < rect.Max.X-n; i++ {
		if i+n < rect.Max.X {
			// We don't need to clone c here because we're just moving cells to
			// the left.
			// b.lines[y][i] = b.lines[y][i+n]
			b.setCell(i, y, b.lines[y][i+n], false)
		}
	}

	// Fill the vacated positions with the given cell
	for i := rect.Max.X - n; i < rect.Max.X; i++ {
		b.setCell(i, y, c, true)
	}
}
