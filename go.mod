module github.com/risor-io/risor

go 1.20

require github.com/stretchr/testify v1.8.3

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.1 // ignores Tamarin release
	v1.0.0 // ignores Tamarin release
)

replace (
	github.com/risor-io/risor/cmd/risor => ./cmd/risor
	github.com/risor-io/risor/modules/aws => ./modules/aws
	github.com/risor-io/risor/modules/image => ./modules/image
	github.com/risor-io/risor/modules/pgx => ./modules/pgx
	github.com/risor-io/risor/modules/uuid => ./modules/uuid
	github.com/risor-io/risor/os/s3fs => ./os/s3fs
)
