module github.com/risor-io/risor/modules/pgx

go 1.20

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.4.1
	github.com/risor-io/risor v0.14.1-0.20230825185206-8956c356a975
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
