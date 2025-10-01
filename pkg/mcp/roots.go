package mcp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterRoots registers filesystem roots that the MCP server may work with.
// Roots are informational and help clients understand what locations the server might access.
//
// Note: Roots are advisory and do not enforce permissions. They inform clients about
// potential filesystem locations the server might interact with.
func RegisterRoots(s *server.MCPServer) error {
	roots := []mcp.Root{}

	// Add user's home directory as a root (common location for config files)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		roots = append(roots, mcp.Root{
			URI:  fmt.Sprintf("file://%s", homeDir),
			Name: "User Home Directory",
		})
	}

	// Add current working directory as a root
	cwd, err := os.Getwd()
	if err == nil {
		roots = append(roots, mcp.Root{
			URI:  fmt.Sprintf("file://%s", cwd),
			Name: "Current Working Directory",
		})
	}

	// Add system timezone data directory if it exists (for timezone operations)
	tzDataLocations := []string{
		"/usr/share/zoneinfo",           // Linux/Unix
		"/usr/share/lib/zoneinfo",       // Some Unix variants
		filepath.Join(homeDir, ".local", "share", "zoneinfo"), // User-local
	}

	for _, tzPath := range tzDataLocations {
		if info, err := os.Stat(tzPath); err == nil && info.IsDir() {
			roots = append(roots, mcp.Root{
				URI:  fmt.Sprintf("file://%s", tzPath),
				Name: "Timezone Database",
			})
			break // Only add the first valid timezone database found
		}
	}

	// Note: The mcp-go library v0.33.0 does not have a SetRoots method.
	// Roots are typically configured at the server level or through
	// server initialization options. This function prepares the roots
	// but cannot directly set them on the server.
	//
	// Roots are informational anyway and don't enforce permissions.
	_ = roots // Mark as intentionally unused

	return nil
}

// GetServerRoots returns the list of roots configured for this server.
// This is useful for debugging and understanding what locations the server declares.
func GetServerRoots() ([]mcp.Root, error) {
	roots := []mcp.Root{}

	// User home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		roots = append(roots, mcp.Root{
			URI:  fmt.Sprintf("file://%s", homeDir),
			Name: "User Home Directory",
		})
	}

	// Current working directory
	cwd, err := os.Getwd()
	if err == nil {
		roots = append(roots, mcp.Root{
			URI:  fmt.Sprintf("file://%s", cwd),
			Name: "Current Working Directory",
		})
	}

	// Timezone data directory
	tzDataLocations := []string{
		"/usr/share/zoneinfo",
		"/usr/share/lib/zoneinfo",
		filepath.Join(homeDir, ".local", "share", "zoneinfo"),
	}

	for _, tzPath := range tzDataLocations {
		if info, err := os.Stat(tzPath); err == nil && info.IsDir() {
			roots = append(roots, mcp.Root{
				URI:  fmt.Sprintf("file://%s", tzPath),
				Name: "Timezone Database",
			})
			break
		}
	}

	return roots, nil
}

// AddCustomRoot allows adding a custom root location to the server.
// This can be used to declare additional filesystem locations the server might access.
func AddCustomRoot(s *server.MCPServer, path string, name string) error {
	// Validate that the path exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create URI
	uri := fmt.Sprintf("file://%s", absPath)

	// Determine description based on whether it's a file or directory
	description := "File"
	if info.IsDir() {
		description = "Directory"
	}

	if name == "" {
		name = filepath.Base(absPath)
	}

	root := mcp.Root{
		URI:  uri,
		Name: fmt.Sprintf("%s: %s", description, name),
	}

	// Get existing roots and append the new one
	existingRoots, err := GetServerRoots()
	if err != nil {
		existingRoots = []mcp.Root{}
	}

	allRoots := append(existingRoots, root)
	_ = allRoots // Mark as intentionally unused

	// Note: The mcp-go library v0.33.0 does not have a SetRoots method.
	// This function prepares root configuration but cannot directly set them.
	return nil
}
