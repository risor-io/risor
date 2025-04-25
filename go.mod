module github.com/risor-io/risor

go 1.23.0

require (
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.3.1
	github.com/fatih/color v1.18.0
	github.com/mattn/go-isatty v0.0.20
	github.com/olekukonko/tablewriter v0.0.5
	github.com/risor-io/risor/modules/gha v1.7.0
	github.com/stretchr/testify v1.10.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/JohannesKaufmann/dom v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

retract (
	v1.0.1 // ignores Tamarin release
	v1.0.0 // ignores Tamarin release
)
