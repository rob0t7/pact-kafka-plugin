VERSION?=0.1.0
PROJECT=kafka
.DEFAULT_GOAL := ci

ci:: install-tools clean build lint test

install-tools:
	@echo "Install local development tools"
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/hamba/avro/v2/cmd/avrosv@latest
	go install github.com/hamba/avro/v2/cmd/avrogen@latest

.PHONY: avrogen
avrogen: install-tools
	avrogen -pkg test -o test/avro_gen.go -tags json:camel,yaml:camel test/schema.avsc

.PHONY: build
build:
	go build -o build/$(PROJECT)

clean:
	rm -rf build dist

proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/pact_plugin.proto

test: install_local
	gotestsum ./...

.PHONY: test-coverage
test-coverage:
	gotestsum -- -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

install_local: build write_config
	@echo "Creating a local phony plugin install in order to test locally"
	mkdir -p ~/.pact/plugins/$(PROJECT)-$(VERSION)/
	cp ./build/$(PROJECT) ~/.pact/plugins/$(PROJECT)-$(VERSION)/
	cp pact-plugin.json ~/.pact/plugins/$(PROJECT)-$(VERSION)/

write_config:
	@cp pact-plugin.json pact-plugin.json.new
	@cat pact-plugin.json | jq '.version = "'$(subst v,,$(VERSION))'" | .name = "'$(PROJECT)'" | .entryPoint = "'$(PROJECT)'"' | tee pact-plugin.json.new
	@mv pact-plugin.json.new pact-plugin.json

.PHONY: build test clean write_config lint
