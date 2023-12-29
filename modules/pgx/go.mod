module github.com/risor-io/risor/modules/pgx

go 1.21

replace github.com/risor-io/risor => ../..

require (
	github.com/jackc/pgx/v5 v5.5.0
	github.com/risor-io/risor v1.1.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
