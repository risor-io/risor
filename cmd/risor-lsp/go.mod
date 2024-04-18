module github.com/risor-io/risor/cmd/risor-lsp

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jdbaldry/go-language-server-protocol v0.0.0-20211013214444-3022da0884b2
	github.com/risor-io/risor v1.5.0
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
)
