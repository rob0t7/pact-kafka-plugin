package main

import (
	"log"
	"log/slog"
	"os"
)

const (
	AVRO_SCHEMA_CONTENT_TYPE = "application/vnd.kafka.avro.v2"
	PLUGIN_NAME              = "kafkaplugin"
)

func main() {
	if err := InitLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	if err := StartPluginServer(); err != nil {
		slog.Error("failed to start plugin server", "error", err)
		os.Exit(1)
	}
}
