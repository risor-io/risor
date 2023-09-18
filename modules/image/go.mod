module github.com/risor-io/risor/modules/image

go 1.21

toolchain go1.21.0

replace github.com/risor-io/risor => ../..

require (
	github.com/anthonynsimon/bild v0.13.0
	github.com/risor-io/risor v1.1.0
)

require golang.org/x/image v0.5.0 // indirect
