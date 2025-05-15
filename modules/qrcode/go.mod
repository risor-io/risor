module github.com/risor-io/risor/modules/qrcode

go 1.23.0

replace (
	github.com/risor-io/risor => ../..
	github.com/risor-io/risor/modules/image => ../image
)

require (
	github.com/risor-io/risor v1.7.0
	github.com/risor-io/risor/modules/image v0.0.0-00010101000000-000000000000
	github.com/yeqown/go-qrcode/v2 v2.2.5
	github.com/yeqown/go-qrcode/writer/standard v1.3.0
)

require (
	github.com/anthonynsimon/bild v0.14.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.27.0 // indirect
)
