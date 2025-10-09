package main

import (
	"log"
	"log/slog"
	"os"
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
