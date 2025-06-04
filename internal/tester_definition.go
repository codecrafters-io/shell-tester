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
		{
			Slug:     "qp2",
			TestFunc: testA1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "gm9",
			TestFunc: testA2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "qm8",
			TestFunc: testA3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "gy5",
			TestFunc: testA4,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "wh6",
			TestFunc: testA5,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "wt6",
			TestFunc: testA6,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "br6",
			TestFunc: testP1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ny9",
			TestFunc: testP2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "xk3",
			TestFunc: testP3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "bq4",
			TestFunc: testH1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "yf5",
			TestFunc: testH2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "ag6",
			TestFunc: testH3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "rh7",
			TestFunc: testH4,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "vq0",
			TestFunc: testH5,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "dm2",
			TestFunc: testH6,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "za2",
			TestFunc: testHP1,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "in3",
			TestFunc: testHP2,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "sx3",
			TestFunc: testHP3,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "zp4",
			TestFunc: testHP4,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "kz7",
			TestFunc: testHP5,
			Timeout:  15 * time.Second,
		},
		{
			Slug:     "jv2",
			TestFunc: testHP6,
			Timeout:  15 * time.Second,
		},
	},
}
