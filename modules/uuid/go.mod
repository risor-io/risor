module github.com/risor-io/risor/modules/uuid

go 1.21

replace github.com/risor-io/risor => ../..

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/risor-io/risor v1.1.0
)

require github.com/russross/blackfriday/v2 v2.1.0 // indirect
