module github.com/risor-io/risor/modules/pgx

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.7.4
	github.com/risor-io/risor v1.7.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
