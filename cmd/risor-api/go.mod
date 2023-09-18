module github.com/risor-io/risor/cmd/risor-api

go 1.21

toolchain go1.21.0

replace github.com/risor-io/risor => ../..

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/risor-io/risor v1.1.0
	github.com/stretchr/testify v1.8.3
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
