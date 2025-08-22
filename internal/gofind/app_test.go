package gofind

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParsePath(t *testing.T) {
	// Get the actual home directory for comparison
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		checkFunc func(t *testing.T, result string)
	}{
		{
			name:    "absolute path unchanged",
			input:   "/usr/local/bin",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				if result != "/usr/local/bin" {
					t.Errorf("Expected /usr/local/bin, got %q", result)
				}
			},
		},
		{
			name:    "relative path unchanged",
			input:   "relative/path",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				if result != "relative/path" {
					t.Errorf("Expected relative/path, got %q", result)
				}
			},
		},
		{
			name:    "tilde expansion - home only",
			input:   "~",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				expected := homeDir
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name:    "tilde expansion - home with subpath",
			input:   "~/Documents",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				expected := filepath.Join(homeDir, "Documents")
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name:    "tilde expansion - nested path",
			input:   "~/Documents/Projects/go",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				expected := filepath.Join(homeDir, "Documents", "Projects", "go")
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("Expected empty string, got %q", result)
				}
			},
		},
		{
			name:    "just tilde with slash",
			input:   "~/",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				expected := filepath.Join(homeDir, "")
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
	}

	// Add Windows-specific test
	if runtime.GOOS == "windows" {
		tests = append(tests, struct {
			name      string
			input     string
			wantErr   bool
			checkFunc func(t *testing.T, result string)
		}{
			name:    "windows absolute path",
			input:   "C:\\Program Files\\Test",
			wantErr: false,
			checkFunc: func(t *testing.T, result string) {
				if result != "C:\\Program Files\\Test" {
					t.Errorf("Expected C:\\Program Files\\Test, got %q", result)
				}
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePath(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			tt.checkFunc(t, result)
		})
	}
}

func TestParsePathTildeEdgeCases(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde not at start",
			input:    "/path/~/notexpanded",
			expected: "/path/~/notexpanded", // Should not be expanded
		},
		{
			name:     "multiple tildes",
			input:    "~/path/~/should/not/expand",
			expected: filepath.Join(homeDir, "path", "~", "should", "not", "expand"), // Only first ~ should be considered
		},
		{
			name:     "tilde with no slash",
			input:    "~noexpand",
			expected: filepath.Join(homeDir, "noexpand"), // Current implementation expands any path starting with ~
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePath(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestParsePathErrorHandling tests error scenarios
// Note: This is harder to test since os.UserHomeDir() rarely fails in normal conditions
func TestParsePathErrorHandling(t *testing.T) {
	// We can't easily make os.UserHomeDir() fail, but we can test the function structure
	// This test mainly ensures the function signature is correct and returns appropriate types
	
	result, err := ParsePath("~/test")
	if err != nil {
		t.Errorf("ParsePath with tilde should not fail in normal conditions: %v", err)
	}
	
	if result == "" {
		t.Error("ParsePath should return non-empty result for valid tilde expansion")
	}
	
	// Test that it returns proper types
	if result == "~/test" {
		t.Error("Tilde should have been expanded, but wasn't")
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		match    Match
		expected struct {
			name string
			path string
		}
	}{
		{
			name:  "basic match",
			match: Match{Name: "project1", Path: "/home/user/project1"},
			expected: struct {
				name string
				path string
			}{
				name: "project1",
				path: "/home/user/project1",
			},
		},
		{
			name:  "empty match",
			match: Match{Name: "", Path: ""},
			expected: struct {
				name string
				path string
			}{
				name: "",
				path: "",
			},
		},
		{
			name:  "match with special characters",
			match: Match{Name: "my-project_v2", Path: "/home/user/projects/my-project_v2"},
			expected: struct {
				name string
				path string
			}{
				name: "my-project_v2",
				path: "/home/user/projects/my-project_v2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.match.Name != tt.expected.name {
				t.Errorf("Expected Name %q, got %q", tt.expected.name, tt.match.Name)
			}
			
			if tt.match.Path != tt.expected.path {
				t.Errorf("Expected Path %q, got %q", tt.expected.path, tt.match.Path)
			}
		})
	}
}

func TestMatchSliceOperations(t *testing.T) {
	matches := []Match{
		{Name: "project-a", Path: "/home/user/project-a"},
		{Name: "project-b", Path: "/home/user/project-b"},
		{Name: "project-c", Path: "/home/user/project-c"},
	}

	// Test that we can iterate over matches
	count := 0
	for _, match := range matches {
		count++
		if match.Name == "" {
			t.Error("Match should have non-empty name")
		}
		if match.Path == "" {
			t.Error("Match should have non-empty path")
		}
	}

	if count != 3 {
		t.Errorf("Expected to iterate over 3 matches, got %d", count)
	}

	// Test that we can append matches
	newMatch := Match{Name: "project-d", Path: "/home/user/project-d"}
	matches = append(matches, newMatch)

	if len(matches) != 4 {
		t.Errorf("Expected 4 matches after append, got %d", len(matches))
	}

	// Test that the appended match is correct
	lastMatch := matches[len(matches)-1]
	if lastMatch.Name != "project-d" {
		t.Errorf("Expected last match name to be 'project-d', got %q", lastMatch.Name)
	}
}