.PHONY: install-tools
install-tools:
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: build
build:
	go build -o pact-kafka-plugin .

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
