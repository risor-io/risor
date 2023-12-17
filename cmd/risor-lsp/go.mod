module github.com/risor-io/risor/cmd/risor-lsp

go 1.21

replace github.com/risor-io/risor => ../..

require (
	github.com/jdbaldry/go-language-server-protocol v0.0.0-20211013214444-3022da0884b2
	github.com/risor-io/risor v1.1.0
	github.com/rs/zerolog v1.30.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
)
