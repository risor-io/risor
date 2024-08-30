module github.com/risor-io/risor/modules/pgx

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.6.0
	github.com/risor-io/risor v1.6.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)
