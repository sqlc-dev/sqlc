package sqlfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	testdataDir := "testdata"

	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testName := entry.Name()
		t.Run(testName, func(t *testing.T) {
			testDir := filepath.Join(testdataDir, testName)

			// Read input file
			inputPath := filepath.Join(testDir, "input.sql")
			inputData, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			// Read expected output files
			var expected []string
			for i := 1; ; i++ {
				outputPath := filepath.Join(testDir, fmt.Sprintf("output_%d.sql", i))
				data, err := os.ReadFile(outputPath)
				if err != nil {
					if os.IsNotExist(err) {
						break
					}
					t.Fatalf("Failed to read output file %s: %v", outputPath, err)
				}
				expected = append(expected, string(data))
			}

			// Run Split
			ctx := context.Background()
			reader := strings.NewReader(string(inputData))

			got, err := Split(ctx, reader)
			if err != nil {
				t.Fatalf("Split() error = %v", err)
			}

			// Compare results
			if len(got) != len(expected) {
				t.Errorf("Split() got %d queries, expected %d", len(got), len(expected))
				t.Logf("Got: %v", got)
				t.Logf("Expected: %v", expected)
				return
			}

			for i := range got {
				if got[i] != expected[i] {
					t.Errorf("Query %d:\ngot:      %q\nexpected: %q", i, got[i], expected[i])
				}
			}
		})
	}
}

func TestSplitContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	reader := strings.NewReader("SELECT * FROM users;")
	_, err := Split(ctx, reader)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestExtractDollarTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty dollar quote",
			input:    "$$",
			expected: "$$",
		},
		{
			name:     "simple tag",
			input:    "$tag$",
			expected: "$tag$",
		},
		{
			name:     "tag with numbers",
			input:    "$tag123$",
			expected: "$tag123$",
		},
		{
			name:     "tag with underscore",
			input:    "$my_tag$",
			expected: "$my_tag$",
		},
		{
			name:     "not a dollar quote (no closing)",
			input:    "$tag",
			expected: "",
		},
		{
			name:     "not a dollar quote (invalid char)",
			input:    "$tag-name$",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no dollar sign",
			input:    "tag",
			expected: "",
		},
		{
			name:     "tag with extra content",
			input:    "$tag$rest of string",
			expected: "$tag$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDollarTag(tt.input)
			if got != tt.expected {
				t.Errorf("extractDollarTag(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}
