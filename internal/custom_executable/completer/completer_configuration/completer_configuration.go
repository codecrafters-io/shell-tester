package completer_configuration

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
	OutputLines []string
	// UseStderrStream will use os.Stderr to print the output lines
	// Using stderr stream will also make the script sleep for 120s
	// to make sure that the error from stderr is streamed to the shell and is
	// not collected after it's exitted instead
	UseStderrStream   bool
	ExpectedArguments *CompleterConfigurationExpectedArguments
	ExpectedEnvVars   *CompleterConfigurationEnvVars
}
