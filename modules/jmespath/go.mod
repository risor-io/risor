module github.com/risor-io/risor/modules/jmespath

go 1.21

replace github.com/risor-io/risor => ../..

require (
	github.com/jmespath-community/go-jmespath v1.1.1
	github.com/risor-io/risor v0.0.0-00010101000000-000000000000
)

require golang.org/x/exp v0.0.0-20230314191032-db074128a8ec // indirect
