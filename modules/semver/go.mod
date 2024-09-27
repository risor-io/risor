module github.com/risor-io/risor/modules/semver

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/risor-io/risor v1.6.0
)
