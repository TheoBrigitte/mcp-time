package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetServerRoots(t *testing.T) {
	roots, err := GetServerRoots()
	if err != nil {
		t.Fatalf("GetServerRoots() error = %v", err)
	}

	if len(roots) == 0 {
		t.Error("GetServerRoots() returned empty roots list")
	}

	// Verify that roots contain expected entries
	foundHome := false
	foundCwd := false

	for _, root := range roots {
		if root.URI == "" {
			t.Error("GetServerRoots() returned root with empty URI")
		}
		if root.Name == "" {
			t.Error("GetServerRoots() returned root with empty Name")
		}

		// Check for home directory
		if root.Name == "User Home Directory" {
			foundHome = true
			homeDir, _ := os.UserHomeDir()
			expectedURI := "file://" + homeDir
			if root.URI != expectedURI {
				t.Errorf("Home directory root URI = %v, want %v", root.URI, expectedURI)
			}
		}

		// Check for current working directory
		if root.Name == "Current Working Directory" {
			foundCwd = true
			cwd, _ := os.Getwd()
			expectedURI := "file://" + cwd
			if root.URI != expectedURI {
				t.Errorf("CWD root URI = %v, want %v", root.URI, expectedURI)
			}
		}
	}

	if !foundHome {
		t.Error("GetServerRoots() did not return User Home Directory root")
	}

	if !foundCwd {
		t.Error("GetServerRoots() did not return Current Working Directory root")
	}
}

func TestRegisterRoots(t *testing.T) {
	// Test that RegisterRoots can be called without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterRoots() panicked: %v", r)
		}
	}()

	s := NewServer("test-server", "1.0.0")
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// RegisterRoots is already called in NewServer
	// Call it again to test it can be called multiple times
	err := RegisterRoots(s.MCPServer)
	if err != nil {
		t.Errorf("RegisterRoots() error = %v", err)
	}
}

func TestAddCustomRoot(t *testing.T) {
	s := NewServer("test-server", "1.0.0")
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	tests := []struct {
		name      string
		path      string
		rootName  string
		wantError bool
		setup     func() (string, func())
	}{
		{
			name:      "add valid directory",
			path:      "",
			rootName:  "Test Directory",
			wantError: false,
			setup: func() (string, func()) {
				tmpDir, err := os.MkdirTemp("", "mcp-test-*")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				cleanup := func() { os.RemoveAll(tmpDir) }
				return tmpDir, cleanup
			},
		},
		{
			name:      "add valid file",
			path:      "",
			rootName:  "Test File",
			wantError: false,
			setup: func() (string, func()) {
				tmpFile, err := os.CreateTemp("", "mcp-test-*.txt")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				path := tmpFile.Name()
				tmpFile.Close()
				cleanup := func() { os.Remove(path) }
				return path, cleanup
			},
		},
		{
			name:      "add non-existent path",
			path:      "/non/existent/path/that/does/not/exist",
			rootName:  "Invalid Path",
			wantError: true,
			setup: func() (string, func()) {
				return "/non/existent/path/that/does/not/exist", func() {}
			},
		},
		{
			name:      "add current directory",
			path:      ".",
			rootName:  "Current Dir",
			wantError: false,
			setup: func() (string, func()) {
				return ".", func() {}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setup()
			defer cleanup()

			if tt.path != "" {
				path = tt.path
			}

			err := AddCustomRoot(s.MCPServer, path, tt.rootName)

			if tt.wantError {
				if err == nil {
					t.Error("AddCustomRoot() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("AddCustomRoot() unexpected error = %v", err)
			}
		})
	}
}

func TestRootURIFormat(t *testing.T) {
	roots, err := GetServerRoots()
	if err != nil {
		t.Fatalf("GetServerRoots() error = %v", err)
	}

	for _, root := range roots {
		// All URIs should start with "file://"
		if len(root.URI) < 7 || root.URI[:7] != "file://" {
			t.Errorf("Root URI %q does not start with 'file://'", root.URI)
		}

		// Extract path from URI
		path := root.URI[7:]

		// Path should be absolute
		if !filepath.IsAbs(path) {
			t.Errorf("Root path %q is not absolute", path)
		}
	}
}

func TestTimezoneDataRoots(t *testing.T) {
	// Test that timezone data directories are included if they exist
	roots, err := GetServerRoots()
	if err != nil {
		t.Fatalf("GetServerRoots() error = %v", err)
	}

	// Check if any timezone database root exists
	foundTzData := false
	for _, root := range roots {
		if root.Name == "Timezone Database" {
			foundTzData = true

			// Verify the path exists
			path := root.URI[7:] // Remove "file://" prefix
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("Timezone Database path %q does not exist: %v", path, err)
			} else if !info.IsDir() {
				t.Errorf("Timezone Database path %q is not a directory", path)
			}
			break
		}
	}

	// It's okay if no timezone database is found on this system
	if !foundTzData {
		t.Log("No timezone database root found (this is okay on some systems)")
	}
}

func TestRootInformation(t *testing.T) {
	// Test that roots provide useful information
	roots, err := GetServerRoots()
	if err != nil {
		t.Fatalf("GetServerRoots() error = %v", err)
	}

	for _, root := range roots {
		t.Logf("Root: %s - %s", root.Name, root.URI)

		// Verify the root has required fields
		if root.URI == "" {
			t.Error("Root has empty URI")
		}
		if root.Name == "" {
			t.Error("Root has empty Name")
		}
	}
}

func TestRootType(t *testing.T) {
	// Verify Root type structure
	var root mcp.Root

	root = mcp.Root{
		URI:  "file:///test/path",
		Name: "Test Root",
	}

	if root.URI != "file:///test/path" {
		t.Errorf("Root URI = %v, want file:///test/path", root.URI)
	}

	if root.Name != "Test Root" {
		t.Errorf("Root Name = %v, want Test Root", root.Name)
	}
}
