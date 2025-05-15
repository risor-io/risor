module github.com/risor-io/risor/modules/yaml

go 1.23.0

replace github.com/risor-io/risor => ../..

require (
	github.com/risor-io/risor v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/kr/text v0.2.0 // indirect
