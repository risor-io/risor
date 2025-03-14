module github.com/risor-io/risor/modules/pgx

go 1.23.0

toolchain go1.24.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.7.2
	github.com/risor-io/risor v1.7.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
