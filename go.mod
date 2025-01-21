module github.com/codecrafters-io/shell-tester

go 1.23

toolchain go1.23.4

require (
	github.com/charmbracelet/x/vt v0.0.0-20250117142827-dd310ffa7553
	github.com/codecrafters-io/tester-utils v0.2.40
	github.com/creack/pty v1.1.24
	github.com/fatih/color v1.18.0
	go.chromium.org/luci v0.0.0-20250120041927-f10777ad7454
)

// Use this to test the local version of tester-utils
// replace github.com/codecrafters-io/tester-utils v0.2.22 => /Users/rohitpaulk/experiments/codecrafters/tester-utils

require (
	github.com/charmbracelet/x/ansi v0.7.0 // indirect
	github.com/charmbracelet/x/wcwidth v0.0.0-20250117142827-dd310ffa7553 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Selective vendoring is not supported by go mod
replace github.com/charmbracelet/x/vt => ./vendored/github.com/charmbracelet/x/vt