package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

// InitLogger initializes the default slog logger to output to log/plugin.log
func InitLogger() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create log directory if it doesn't exist
	logDir := filepath.Join(cwd, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Create or open log file
	logFile := filepath.Join(logDir, "plugin.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Create JSON handler that writes to file
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Set as default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
