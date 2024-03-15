module github.com/risor-io/risor/cmd/risor-api

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/risor-io/risor v1.5.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
