package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	ExecutableFileName:       "your_program.sh",
	LegacyExecutableFileName: "your_shell.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "oo8",
			TestFunc: testPrompt,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "cz2",
			TestFunc: testInvalidCommand,
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
		{
			Slug:     "ni6",
			TestFunc: testQ1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "tg6",
			TestFunc: testQ2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "yt5",
			TestFunc: testQ3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "le5",
			TestFunc: testQ4,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "gu3",
			TestFunc: testQ5,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "qj0",
			TestFunc: testQ6,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "jv1",
			TestFunc: testR1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "vz4",
			TestFunc: testR2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "el9",
			TestFunc: testR3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "un3",
			TestFunc: testR4,
			Timeout:  15 * time.Second,
		},
	},
}
