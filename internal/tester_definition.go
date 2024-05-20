package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	ExecutableFileName: "spawn_shell.sh",
	TestCases: []tester_definition.TestCase{
		// {
		// 	Slug:     "init",
		// 	TestFunc: testPrompt,
		// 	Timeout:  15 * time.Second,
		// },
		{
			Slug:     "missing-command",
			TestFunc: testMissingCommand,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "repl",
			TestFunc: testREPL,
			Timeout:  15 * time.Second,
		},
		// {
		// 	Slug:     "exit",
		// 	TestFunc: testExit,
		// 	Timeout:  15 * time.Second,
		// },
		// {
		// 	Slug:     "echo",
		// 	TestFunc: testEcho,
		// 	Timeout:  15 * time.Second,
		// },
	},
}
