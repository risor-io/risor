module github.com/risor-io/risor/cmd/risor-lsp

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jdbaldry/go-language-server-protocol v0.0.0-20211013214444-3022da0884b2
	github.com/risor-io/risor v1.8.0
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
)
