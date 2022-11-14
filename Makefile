
.PHONY: bench
bench:
	go build
	./tamarin -profile cpu.out ./benchmark/main.mon
	go tool pprof -http=:8080 ./cpu.out

# https://code.visualstudio.com/api/working-with-extensions/publishing-extension#packaging-extensions
.PHONY: install-tools
install-tools:
	npm install -g vsce

.PHONY: extension-login
extension-login:
	cd vscode && vsce login CurtisMyzie

.PHONY: extension
extension:
	cd vscode && vsce package && vsce publish

.PHONY: postgres
postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
