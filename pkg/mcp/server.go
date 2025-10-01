package mcp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/mark3labs/mcp-go/server"
	"github.com/mark3labs/mcp-go/util"
)

// transport defines the communication protocol for the MCP server.
type transport int

const (
	// TransportSTDIO uses standard input/output for communication.
	TransportSTDIO transport = iota
	// TransportStream uses an HTTP stream for communication.
	TransportStream
)

// TransportNames maps transport types to their string representations.
var TransportNames = map[transport]string{
	TransportSTDIO:  "stdio",
	TransportStream: "stream",
}

// GetTransports returns a slice of available transport names.
func GetTransports() []string { return slices.Collect(maps.Values(TransportNames)) }

// Server wraps the core MCP server and registers the time-specific handlers.
type Server struct {
	*server.MCPServer
}

// NewServer creates a new MCP server with the time tools registered.
// It initializes the underlying MCP server and registers all the tool handlers.
func NewServer(name, version string) *Server {
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, false), // resources enabled, subscriptions disabled
	)

	RegisterHandlers(mcpServer)

	// Register filesystem roots
	err := RegisterRoots(mcpServer)
	if err != nil {
		log.Printf("Warning: failed to register roots: %v", err)
	}

	s := &Server{
		mcpServer,
	}

	return s
}

// StartStdio starts the server listening on standard input/output.
// It blocks until the context is canceled or an error occurs.
func (s Server) StartStdio(ctx context.Context) error {
	stdioServer := server.NewStdioServer(s.MCPServer)

	// Set the logger for the stdio server.
	stdioServer.SetErrorLogger(log.Default())

	err := stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	// Don't return an error if the context was canceled, as it's an expected shutdown.
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

// StartStream starts the server as an HTTP stream on the given address.
// It includes a graceful shutdown mechanism that is triggered by the context.
func (s Server) StartStream(ctx context.Context, address string) error {
	u, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("invalid address %s: %w", address, err)
	}

	if u.Port() == "" {
		return fmt.Errorf("invalid address %s: expected format scheme://host:port", address)
	}

	hostPort := net.JoinHostPort(u.Hostname(), u.Port())

	streamServer := server.NewStreamableHTTPServer(s.MCPServer,
		server.WithLogger(util.DefaultLogger()),
		server.WithEndpointPath(u.Path),
	)

	go func() {
		<-ctx.Done()
		// Gracefully shut down the server when the context is canceled.
		streamServer.Shutdown(ctx) // nolint:gosec,errcheck
	}()

	err = streamServer.Start(hostPort)
	// Don't return an error on a clean server shutdown.
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
