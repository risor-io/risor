module github.com/risor-io/risor/modules/bcrypt

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/risor-io/risor v1.5.2
	golang.org/x/crypto v0.22.0
)
