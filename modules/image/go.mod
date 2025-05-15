module github.com/risor-io/risor/modules/image

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/anthonynsimon/bild v0.14.0
	github.com/risor-io/risor v1.8.0
)

require golang.org/x/image v0.27.0 // indirect
