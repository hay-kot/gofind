package ui

import (
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected struct {
			hasNewlines bool
			hasSpaces   bool
			rowCount    int
		}
	}{
		{
			name: "simple 2x2 table",
			input: [][]string{
				{"Name", "Age"},
				{"John", "25"},
			},
			expected: struct {
				hasNewlines bool
				hasSpaces   bool
				rowCount    int
			}{
				hasNewlines: true,
				hasSpaces:   true,
				rowCount:    2,
			},
		},
		{
			name: "table with varying column widths",
			input: [][]string{
				{"Short", "Very Long Column Name", "Mid"},
				{"A", "B", "C"},
				{"123", "456789", "XY"},
			},
			expected: struct {
				hasNewlines bool
				hasSpaces   bool
				rowCount    int
			}{
				hasNewlines: true,
				hasSpaces:   true,
				rowCount:    3,
			},
		},
		{
			name: "single row table (header only)",
			input: [][]string{
				{"Header1", "Header2", "Header3"},
			},
			expected: struct {
				hasNewlines bool
				hasSpaces   bool
				rowCount    int
			}{
				hasNewlines: true,
				hasSpaces:   true,
				rowCount:    1,
			},
		},
		{
			name: "table with empty strings",
			input: [][]string{
				{"Name", "", "Email"},
				{"", "Doe", "john@example.com"},
			},
			expected: struct {
				hasNewlines bool
				hasSpaces   bool
				rowCount    int
			}{
				hasNewlines: true,
				hasSpaces:   true,
				rowCount:    2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Table(tt.input)
			
			if result == "" {
				t.Error("Table() returned empty string")
			}
			
			// Check for newlines
			if tt.expected.hasNewlines && !strings.Contains(result, "\n") {
				t.Error("Expected result to contain newlines")
			}
			
			// Check row count by counting newlines
			newlineCount := strings.Count(result, "\n")
			if newlineCount != tt.expected.rowCount {
				t.Errorf("Expected %d newlines (rows), got %d", tt.expected.rowCount, newlineCount)
			}
			
			// Check that content includes original data
			for _, row := range tt.input {
				for _, cell := range row {
					if cell != "" && !strings.Contains(result, cell) {
						t.Errorf("Expected result to contain %q", cell)
					}
				}
			}
		})
	}
}

func TestTableAlignment(t *testing.T) {
	input := [][]string{
		{"A", "BB", "CCC"},
		{"1234", "5", "67"},
	}
	
	result := Table(input)
	lines := strings.Split(result, "\n")
	
	// Remove empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	
	if len(nonEmptyLines) != 2 {
		t.Fatalf("Expected 2 non-empty lines, got %d", len(nonEmptyLines))
	}
	
	// All lines should have roughly similar lengths due to padding
	firstLineLen := len(nonEmptyLines[0])
	secondLineLen := len(nonEmptyLines[1])
	
	// Allow some variance but they should be similar
	if abs(firstLineLen-secondLineLen) > 10 {
		t.Errorf("Lines have very different lengths: %d vs %d", firstLineLen, secondLineLen)
	}
}

func TestTablePadding(t *testing.T) {
	input := [][]string{
		{"Short", "VeryLongText"},
		{"X", "Y"},
	}
	
	result := Table(input)
	
	// The result should contain appropriate spacing
	// "Short" should be padded to at least match "VeryLongText" length + 4 spaces
	if !strings.Contains(result, "Short    ") { // At least 4 extra spaces
		t.Error("Expected proper padding for shorter columns")
	}
	
	// Should contain all original text
	if !strings.Contains(result, "VeryLongText") {
		t.Error("Expected to contain original long text")
	}
	
	if !strings.Contains(result, "Short") {
		t.Error("Expected to contain original short text")
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestTableEdgeCases(t *testing.T) {
	t.Run("single cell", func(t *testing.T) {
		input := [][]string{{"OnlyCell"}}
		result := Table(input)
		
		if !strings.Contains(result, "OnlyCell") {
			t.Error("Expected to contain the single cell content")
		}
		
		if !strings.Contains(result, "\n") {
			t.Error("Expected to end with newline")
		}
	})
	
	t.Run("wide table", func(t *testing.T) {
		input := [][]string{
			{"Col1", "Col2", "Col3", "Col4", "Col5"},
			{"A", "B", "C", "D", "E"},
		}
		result := Table(input)
		
		// Should contain all columns
		for i := 1; i <= 5; i++ {
			expected := "Col" + string(rune('0'+i))
			if !strings.Contains(result, expected) {
				t.Errorf("Expected to contain %q", expected)
			}
		}
	})
}