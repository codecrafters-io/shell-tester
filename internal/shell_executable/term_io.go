package shell_executable

import "os"

// TermIO represents a terminal input/output pair where
// reading occurs from the pseudo-terminal (pty) and
// writing occurs to the virtual terminal (vt)
type TermIO struct {
	vt  *VirtualTerminal
	pty *os.File
}

// TermIO implements the io.Reader interface
// But we want vt and pty to be always in sync, so we write to vt whenever we read from pty

// Read will read from the pty and write to the vt
func (t *TermIO) Read(p []byte) (n int, err error) {
	readBytes, err := t.pty.Read(p)
	if err != nil {
		return readBytes, err
	}

	t.vt.Write(p[:readBytes])

	return readBytes, nil
}

func NewTermIO(vt *VirtualTerminal, pty *os.File) *TermIO {
	return &TermIO{
		vt:  vt,
		pty: pty,
	}
}
