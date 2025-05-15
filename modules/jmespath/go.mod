module github.com/risor-io/risor/modules/jmespath

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jmespath-community/go-jmespath v1.1.1
	github.com/risor-io/risor v1.7.0
)

require golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
