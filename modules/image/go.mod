module github.com/risor-io/risor/modules/image

go 1.21

toolchain go1.21.0

replace github.com/risor-io/risor => ../..

require (
	github.com/anthonynsimon/bild v0.13.0
	github.com/risor-io/risor v0.14.1-0.20230825185206-8956c356a975
)

require golang.org/x/image v0.5.0 // indirect
