module github.com/risor-io/risor/modules/uuid

go 1.21

toolchain go1.21.0

replace github.com/risor-io/risor => ../..

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/risor-io/risor v1.1.0
)
