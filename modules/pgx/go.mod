module github.com/risor-io/risor/modules/pgx

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.5.5
	github.com/risor-io/risor v1.6.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
