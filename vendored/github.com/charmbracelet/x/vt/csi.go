package vt

import (
	"github.com/charmbracelet/x/ansi"
)

// handleCsi handles a CSI escape sequences.
func (t *Terminal) handleCsi(seq ansi.CsiSequence) {
	switch cmd := t.parser.Cmd(); cmd { // cursor
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'Z', 'a', 'b', 'd', 'e', 'f', '`':
		t.handleCursor()
	case 'm': // Select Graphic Rendition [ansi.SGR]
		t.handleSgr()
	case 'J', 'L', 'M', 'X', 'r', 's':
		t.handleScreen()
	case 'K', 'S', 'T':
		t.handleLine()
	case ansi.Cmd(0, 0, 'h'), ansi.Cmd('?', 0, 'h'): // Set Mode [ansi.SM]
		fallthrough
	case ansi.Cmd(0, 0, 'l'), ansi.Cmd('?', 0, 'l'): // Reset Mode [ansi.RM]
		t.handleMode()
	case ansi.Cmd('?', 0, 'W'): // Set Tab at Every 8 Columns [ansi.DECST8C]
		if params := t.parser.Params(); len(params) == 1 && params[0] == 5 {
			t.resetTabStops()
		}
	case ansi.Cmd(0, ' ', 'q'): // Set Cursor Style [ansi.DECSCUSR]
		style := 1
		if param, ok := t.parser.Param(0, 0); ok {
			style = param
		}
		t.scr.setCursorStyle(CursorStyle((style/2)+1), style%2 == 1)
	case 'g': // Tab Clear [ansi.TBC]
		var value int
		if param, ok := t.parser.Param(0, 0); ok {
			value = param
		}

		switch value {
		case 0:
			x, _ := t.scr.CursorPosition()
			t.tabstops.Reset(x)
		case 3:
			t.tabstops.Clear()
		}
	case '@': // Insert Character [ansi.ICH]
		n, _ := t.parser.Param(0, 1)
		if n == 0 {
			n = 1
		}

		t.scr.InsertCell(n)
	case 'P': // Delete Character [ansi.DCH]
		n, _ := t.parser.Param(0, 1)
		if n == 0 {
			n = 1
		}

		t.scr.DeleteCell(n)

	case 'c': // Primary Device Attributes [ansi.DA1]
		n, _ := t.parser.Param(0, 0)
		if n != 0 {
			break
		}

		// Do we fully support VT220?
		t.buf.WriteString(ansi.PrimaryDeviceAttributes(
			62, // VT220
			1,  // 132 columns
			6,  // Selective Erase
			22, // ANSI color
		))

	case ansi.Cmd('>', 0, 'c'): // Secondary Device Attributes [ansi.DA2]
		n, _ := t.parser.Param(0, 0)
		if n != 0 {
			break
		}

		// Do we fully support VT220?
		t.buf.WriteString(ansi.SecondaryDeviceAttributes(
			1,  // VT220
			10, // Version 1.0
			0,  // ROM Cartridge is always zero
		))

	case 'n': // Device Status Report [ansi.DSR]
		n, ok := t.parser.Param(0, 1)
		if !ok || n == 0 {
			break
		}

		switch n {
		case 5: // Operating Status
			// We're always ready ;)
			// See: https://vt100.net/docs/vt510-rm/DSR-OS.html
			t.buf.WriteString(ansi.DeviceStatusReport(ansi.DECStatus(0)))
		case 6: // Cursor Position Report [ansi.CPR]
			x, y := t.scr.CursorPosition()
			t.buf.WriteString(ansi.CursorPositionReport(x+1, y+1))
		}

	case ansi.Cmd('?', 0, 'n'): // Device Status Report (DEC) [ansi.DSR]
		n, ok := t.parser.Param(0, 1)
		if !ok || n == 0 {
			break
		}

		switch n {
		case 6: // Extended Cursor Position Report [ansi.DECXCPR]
			x, y := t.scr.CursorPosition()
			t.buf.WriteString(ansi.ExtendedCursorPositionReport(x+1, y+1, 0)) // We don't support page numbers
		}

	case ansi.Cmd(0, '$', 'p'): // Request Mode [ansi.DECRQM]
		fallthrough
	case ansi.Cmd('?', '$', 'p'): // Request Mode (DEC) [ansi.DECRQM]
		n, ok := t.parser.Param(0, 0)
		if !ok || n == 0 {
			break
		}

		var mode ansi.Mode = ansi.ANSIMode(n)
		if cmd.Marker() == '?' {
			mode = ansi.DECMode(n)
		}

		setting := t.modes[mode]
		t.buf.WriteString(ansi.ReportMode(mode, setting))

	default:
		t.logf("unhandled CSI: %q", seq)
	}
}
