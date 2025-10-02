# pact-kafka-plugin

A Pact plugin for Kafka message verification. It allows for verification of Asynchronous messages pacts using
AVRO serialization with knowledge of the Kafka Schema Registry.

## Compiling Proto Files

To compile the proto files and generate the gRPC stubs:

```bash
make proto
```

This will generate the Go code from the proto files located in the `proto/` directory.

## Prerequisites

### Installing Protobuf Compiler

#### macOS

```bash
# Install protoc using Homebrew
brew install protobuf
```

#### Linux

```bash
# Debian/Ubuntu
sudo apt-get update
sudo apt-get install -y protobuf-compiler

# Fedora/RHEL
sudo dnf install protobuf-compiler

# Or download the latest release from GitHub
# https://github.com/protocolbuffers/protobuf/releases
```

### Installing Go Protobuf Plugins

After installing `protoc`, install the Go plugins (same for both macOS and Linux):

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Ensure `$GOPATH/bin` is in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
