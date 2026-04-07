package completer_configuration

import "fmt"

type CompleterConfigurationExpectedArguments struct {
	Argv1 string
	Argv2 string
	Argv3 string
}

type CompleterConfigurationEnvVars struct {
	CompLine  string
	CompPoint string
}

type CompleterConfiguration struct {
	CompletionCandidates []string
	StderrLines          []string
	ExpectedArguments    *CompleterConfigurationExpectedArguments
	ExpectedEnvVars      *CompleterConfigurationEnvVars
}

func (c *CompleterConfiguration) Verify() error {
	if len(c.CompletionCandidates) > 0 && len(c.StderrLines) > 0 {
		return fmt.Errorf("Codecrafters Internal Error: Completer Configuration cannot have both CompletionCandidates and StderrLines")
	}

	return nil
}
