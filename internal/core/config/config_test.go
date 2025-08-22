package config

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}
	
	if cfg.Commands == nil {
		t.Error("Commands map should be initialized")
	}
	
	if cfg.MaxRecursion != 10 {
		t.Errorf("Expected MaxRecursion to be 10, got %d", cfg.MaxRecursion)
	}
	
	if len(cfg.Ignore) != 0 {
		t.Errorf("Expected empty Ignore slice, got %v", cfg.Ignore)
	}
	
	if cfg.CacheDir == "" {
		t.Error("CacheDir should not be empty")
	}
}

func TestIgnorePatterns(t *testing.T) {
	patterns := IgnorePatterns()
	expected := []string{"node_modules", ".venv", "venv"}
	
	if len(patterns) != len(expected) {
		t.Fatalf("Expected %d patterns, got %d", len(expected), len(patterns))
	}
	
	for i, pattern := range patterns {
		if pattern != expected[i] {
			t.Errorf("Expected pattern %q at index %d, got %q", expected[i], i, pattern)
		}
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected *Config
	}{
		{
			name:  "valid config",
			input: `{"default":"test","commands":{"repo":{"roots":["/home/user"],"match":"*.go"}},"cache":"/tmp/cache","ignore":["node_modules"],"max_recursion":5}`,
			wantErr: false,
			expected: &Config{
				Default: "test",
				Commands: map[string]SearchEntry{
					"repo": {
						Roots:    []string{"/home/user"},
						MatchStr: "*.go",
					},
				},
				CacheDir:     "/tmp/cache",
				Ignore:       []string{"node_modules"},
				MaxRecursion: 5,
			},
		},
		{
			name:    "invalid json",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
		{
			name:  "empty config uses defaults",
			input: `{}`,
			wantErr: false,
			expected: &Config{
				Default:      "",
				Commands:     map[string]SearchEntry{},
				CacheDir:     XDGCachePath(), // This will be the actual XDG path
				Ignore:       []string{},
				MaxRecursion: 10,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			cfg, err := Read(reader)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			// Compare key fields
			if cfg.Default != tt.expected.Default {
				t.Errorf("Expected Default %q, got %q", tt.expected.Default, cfg.Default)
			}
			
			if cfg.MaxRecursion != tt.expected.MaxRecursion {
				t.Errorf("Expected MaxRecursion %d, got %d", tt.expected.MaxRecursion, cfg.MaxRecursion)
			}
			
			if len(cfg.Commands) != len(tt.expected.Commands) {
				t.Errorf("Expected %d commands, got %d", len(tt.expected.Commands), len(cfg.Commands))
			}
			
			for key, expectedEntry := range tt.expected.Commands {
				actualEntry, exists := cfg.Commands[key]
				if !exists {
					t.Errorf("Expected command %q not found", key)
					continue
				}
				
				if actualEntry.MatchStr != expectedEntry.MatchStr {
					t.Errorf("Expected MatchStr %q, got %q", expectedEntry.MatchStr, actualEntry.MatchStr)
				}
				
				if len(actualEntry.Roots) != len(expectedEntry.Roots) {
					t.Errorf("Expected %d roots, got %d", len(expectedEntry.Roots), len(actualEntry.Roots))
				}
			}
		})
	}
}

func TestWrite(t *testing.T) {
	cfg := &Config{
		Default: "test-default",
		Commands: map[string]SearchEntry{
			"repos": {
				Roots:    []string{"/home/user/repos"},
				MatchStr: "*.git",
			},
		},
		CacheDir:     "/tmp/test-cache",
		Ignore:       []string{"node_modules", ".git"},
		MaxRecursion: 15,
	}
	
	var buf bytes.Buffer
	err := Write(&buf, cfg)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	
	// Parse the output to verify it's valid JSON
	var parsed Config
	err = json.Unmarshal(buf.Bytes(), &parsed)
	if err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}
	
	// Verify key fields
	if parsed.Default != cfg.Default {
		t.Errorf("Expected Default %q, got %q", cfg.Default, parsed.Default)
	}
	
	if parsed.MaxRecursion != cfg.MaxRecursion {
		t.Errorf("Expected MaxRecursion %d, got %d", cfg.MaxRecursion, parsed.MaxRecursion)
	}
	
	// Check if output is properly formatted (indented)
	output := buf.String()
	if !strings.Contains(output, "\n") {
		t.Error("Output should be indented with newlines")
	}
}

func TestReadWrite_RoundTrip(t *testing.T) {
	original := &Config{
		Default: "projects",
		Commands: map[string]SearchEntry{
			"repos": {
				Roots:    []string{"/home/user/repos", "/opt/projects"},
				MatchStr: "go.mod",
			},
			"docs": {
				Roots:    []string{"/home/user/documents"},
				MatchStr: "*.md",
			},
		},
		CacheDir:     "/tmp/gofind-cache",
		Ignore:       []string{"node_modules", ".git", "vendor"},
		MaxRecursion: 20,
	}
	
	// Write to buffer
	var buf bytes.Buffer
	err := Write(&buf, original)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	
	// Read back from buffer
	parsed, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	
	// Compare all fields
	if parsed.Default != original.Default {
		t.Errorf("Default mismatch: expected %q, got %q", original.Default, parsed.Default)
	}
	
	if parsed.CacheDir != original.CacheDir {
		t.Errorf("CacheDir mismatch: expected %q, got %q", original.CacheDir, parsed.CacheDir)
	}
	
	if parsed.MaxRecursion != original.MaxRecursion {
		t.Errorf("MaxRecursion mismatch: expected %d, got %d", original.MaxRecursion, parsed.MaxRecursion)
	}
	
	if len(parsed.Ignore) != len(original.Ignore) {
		t.Errorf("Ignore length mismatch: expected %d, got %d", len(original.Ignore), len(parsed.Ignore))
	}
	
	if len(parsed.Commands) != len(original.Commands) {
		t.Errorf("Commands length mismatch: expected %d, got %d", len(original.Commands), len(parsed.Commands))
	}
	
	for key, expectedEntry := range original.Commands {
		actualEntry, exists := parsed.Commands[key]
		if !exists {
			t.Errorf("Command %q missing after round trip", key)
			continue
		}
		
		if actualEntry.MatchStr != expectedEntry.MatchStr {
			t.Errorf("MatchStr mismatch for %q: expected %q, got %q", key, expectedEntry.MatchStr, actualEntry.MatchStr)
		}
		
		if len(actualEntry.Roots) != len(expectedEntry.Roots) {
			t.Errorf("Roots length mismatch for %q: expected %d, got %d", key, len(expectedEntry.Roots), len(actualEntry.Roots))
		}
		
		for i, expectedRoot := range expectedEntry.Roots {
			if i < len(actualEntry.Roots) && actualEntry.Roots[i] != expectedRoot {
				t.Errorf("Root mismatch for %q[%d]: expected %q, got %q", key, i, expectedRoot, actualEntry.Roots[i])
			}
		}
	}
}