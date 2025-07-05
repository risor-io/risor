module github.com/risor-io/risor/modules/redis

go 1.24.0

replace github.com/risor-io/risor => ../..

require (
	github.com/redis/go-redis/v9 v9.8.0
	github.com/risor-io/risor v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
