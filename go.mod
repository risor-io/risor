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
	github.com/risor-io/risor/modules/aws v0.0.0 => ./modules/aws
	github.com/risor-io/risor/modules/image v0.0.0 => ./modules/image
	github.com/risor-io/risor/modules/pgx v0.0.0 => ./modules/pgx
	github.com/risor-io/risor/modules/uuid v0.0.0 => ./modules/uuid
	github.com/risor-io/risor/os/s3fs v0.0.0 => ./os/s3fs
	github.com/risor-io/risor/cmd/risor v0.0.0 => ./cmd/risor
)
