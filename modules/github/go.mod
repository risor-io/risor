module github.com/risor-io/risor/modules/github

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/google/go-github/v73 v73.0.0
	github.com/risor-io/risor v1.8.0
)

require github.com/google/go-querystring v1.1.0 // indirect
