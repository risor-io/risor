module github.com/risor-io/risor/modules/uuid

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/risor-io/risor v1.5.0
)
