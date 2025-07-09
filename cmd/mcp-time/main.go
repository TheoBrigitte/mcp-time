package main

import (
	"log/slog"
	"os"
)

// main is the entry point of the application.
func main() {
	// Execute the root command and handle any errors.
	err := cmd.Execute()
	if err != nil {
		slog.Error("execution error", "error", err)
		os.Exit(1)
	}
}
