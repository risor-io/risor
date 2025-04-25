module github.com/risor-io/risor/modules/qrcode

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/risor-io/risor v0.0.0-00010101000000-000000000000
	github.com/yeqown/go-qrcode/v2 v2.2.5
	github.com/yeqown/go-qrcode/writer/standard v1.2.1
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5 // indirect
)
