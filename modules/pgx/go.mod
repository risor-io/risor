module github.com/risor-io/risor/modules/pgx

go 1.21

toolchain go1.21.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.4.1
	github.com/risor-io/risor v1.1.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
