package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/mcp-time/pkg/mcp"
)

var (
	// Name is the name of the application.
	// It is set at build time using the -X linker flag.
	Name = ""

	// address is the listen address for the HTTP server.
	address string
	// logFile is the path to the log file. If empty, logs are disabled for stdio transport.
	logFile string
	// transport is the transport layer to use for MCP communication.
	transport   string
	versionFlag = false // Flag to enable version output
)

// cmd defines the root command for the MCP time server.
var cmd = &cobra.Command{
	Use:   Name,
	Short: "MCP server providing time and date utilities",
	Long:  `An MCP (Model Context Protocol) server which provides utilities to work with time and dates.`,
	RunE:  runner,
}

// init initializes command line flags for the application.
func init() {
	cmd.Flags().StringVar(&address, "address", "http://localhost:8080/mcp", "Listen address for Stream HTTP Server (only for --transport stream)")
	cmd.Flags().StringVar(&logFile, "log-file", "", "Path to log file (logs is disabled if not specified)")
	cmd.Flags().StringVarP(&transport, "transport", "t", mcp.TransportNames[mcp.TransportSTDIO], fmt.Sprintf("Transport layer: %v.", strings.Join(mcp.GetTransports(), ", ")))
	cmd.Flags().BoolVar(&versionFlag, "version", false, "Print version information and exit")
}

// runner is the main execution function for the MCP server.
// It sets up logging, creates the MCP server, and starts it.
func runner(c *cobra.Command, args []string) (err error) {
	// If the version flag is set, print the version and exit.
	if versionFlag {
		fmt.Println(version.Print(Name))
		return nil
	}

	// Set up logging. Default to a no-op handler.
	var logger = slog.DiscardHandler

	// If a log file is specified, create/open it and use it for logging.
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // nolint:gosec
		if err != nil {
			return err
		}
		defer file.Close() // nolint:errcheck
		logger = slog.NewTextHandler(file, nil)
	} else if transport != mcp.TransportNames[mcp.TransportSTDIO] {
		// For non-stdio transports, log to stderr by default if no log file is provided.
		logger = slog.NewTextHandler(os.Stderr, nil)
	}

	// Set the default logger for the application.
	slog.SetDefault(slog.New(logger))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-sigChan
		slog.Info(fmt.Sprintf("%s received", signal))
		cancel()
	}()

	// Create a new MCP server instance.
	server := mcp.NewServer(Name, version.Version)

	// Start the server with the configured transport.
	switch transport {
	case mcp.TransportNames[mcp.TransportSTDIO]:
		slog.Info("MCP server starting", "version", version.Version, "transport", transport)
		err = server.StartStdio(ctx)
	case mcp.TransportNames[mcp.TransportStream]:
		slog.Info("MCP server starting", "version", version.Version, "transport", transport, "address", address)
		err = server.StartStream(ctx, address)
	default:
		return fmt.Errorf("transport not supported: %s", transport)
	}

	if err != nil {
		return fmt.Errorf("MCP server startup failed: %w", err)
	}

	slog.Info("MCP server shutdown")

	return nil
}
