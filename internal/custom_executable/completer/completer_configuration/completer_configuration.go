package completer_configuration

import "fmt"

type CompleterConfiguration struct {
	CompletionCandidates []string
	StderrLines          []string
}

func (c *CompleterConfiguration) Verify() error {
	if len(c.CompletionCandidates) > 0 && len(c.StderrLines) > 0 {
		return fmt.Errorf("Codecrafters Internal Error: Completer Configuration cannot have both CompletionCandidates and StderrLines")
	}

	return nil
}
