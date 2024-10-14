module github.com/risor-io/risor/modules/image

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/anthonynsimon/bild v0.14.0
	github.com/risor-io/risor v1.7.0
)

require (
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/image v0.19.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
