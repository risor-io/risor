module github.com/risor-io/risor/modules/sql

go 1.22

toolchain go1.22.0

replace github.com/risor-io/risor => ../..

require (
	github.com/go-sql-driver/mysql v1.8.0
	github.com/lib/pq v1.10.9
	github.com/microsoft/go-mssqldb v1.7.0
	github.com/risor-io/risor v1.5.0
	github.com/xo/dburl v0.21.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
