package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	ExecutableFileName: "spawn_shell.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "oo8",
			TestFunc: testPrompt,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "cz2",
			TestFunc: testMissingCommand,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ff0",
			TestFunc: testREPL,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "pn5",
			TestFunc: testExit,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "iz3",
			TestFunc: testEcho,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ez5",
			TestFunc: testType1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "mg5",
			TestFunc: testType2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ip1",
			TestFunc: testRun,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ei0",
			TestFunc: testpwd,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ra6",
			TestFunc: testCd1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "gq9",
			TestFunc: testCd2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "gp4",
			TestFunc: testCd3,
			Timeout:  15 * time.Second,
		},
	},
}
