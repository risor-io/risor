module github.com/risor-io/risor/modules/semver

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/blang/semver/v4 v4.0.0
	github.com/risor-io/risor v1.6.0
)
