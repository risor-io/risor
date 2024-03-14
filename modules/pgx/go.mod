module github.com/risor-io/risor/modules/pgx

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.5.4
	github.com/risor-io/risor v1.1.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
