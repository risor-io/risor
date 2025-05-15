module github.com/risor-io/risor/modules/playwright

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/playwright-community/playwright-go v0.5200.0
	github.com/risor-io/risor v0.0.0-00010101000000-000000000000
)

require (
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
)
