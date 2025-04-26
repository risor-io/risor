module github.com/risor-io/risor/modules/slack

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/risor-io/risor v1.7.0
	github.com/slack-go/slack v0.16.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	golang.org/x/net v0.39.0 // indirect
)
