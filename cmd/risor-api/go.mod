module github.com/risor-io/risor/cmd/risor-api

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/risor-io/risor v1.3.2
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
