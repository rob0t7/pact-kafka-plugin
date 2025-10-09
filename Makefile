.PHONY: install-tools
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
	go build -o build/kafkaplugin

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/pact_plugin.proto

.PHONY: test
test:
	gotestsum ./...

.PHONY: test-coverage
test-coverage:
	gotestsum -- -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint:
	golangci-lint run

.PHONY: install_local
install_local: build
	@echo "Creating a local phony plugin install in order to test locally"
	mkdir -p ~/.pact/plugins/kafka-0.0.1/
	mkdir -p ~/.pact/plugins/kafka-0.0.1/log
	cp ./build/kafkaplugin ~/.pact/plugins/kafka-0.0.1/
	cp pact-plugin.json ~/.pact/plugins/kafka-0.0.1/
