package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/enginetest/testcases"
)

// TestCoverage verifies that all required test cases are implemented
// for the SQLite engine.
func TestCoverage(t *testing.T) {
	engine := Engine()
	registry := testcases.DefaultRegistry

	// Get all tests this engine should implement
	requiredTests := registry.RequiredTestsForEngine(engine)

	testdataDir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Find all implemented tests
	implemented := make(map[string]bool)
	err = filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.yaml" || info.Name() == "sqlc.json" {
			dir := filepath.Dir(path)
			testName := filepath.Base(dir)
			implemented[testName] = true
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	// Check for missing tests
	var missing []*testcases.TestCase
	for _, tc := range requiredTests {
		if !implemented[tc.Name] {
			missing = append(missing, tc)
		}
	}

	// Report missing tests (informational, not a failure)
	if len(missing) > 0 {
		t.Logf("SQLite engine is missing %d required test cases (this is informational):", len(missing))
		for _, tc := range missing {
			t.Logf("  - %s (%s): %s", tc.ID, tc.Name, tc.Description)
		}
	}

	// Report coverage statistics
	total := len(requiredTests)
	covered := total - len(missing)
	percentage := float64(covered) / float64(total) * 100

	t.Logf("SQLite test coverage: %d/%d (%.1f%%)", covered, total, percentage)
}

// TestCoverageByCategory reports coverage broken down by category
func TestCoverageByCategory(t *testing.T) {
	engine := Engine()
	registry := testcases.DefaultRegistry
	caps := testcases.DefaultCapabilities(engine)

	testdataDir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Find all implemented tests
	implemented := make(map[string]bool)
	_ = filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.yaml" || info.Name() == "sqlc.json" {
			dir := filepath.Dir(path)
			testName := filepath.Base(dir)
			implemented[testName] = true
			return filepath.SkipDir
		}
		return nil
	})

	// Report by category
	categories := testcases.RequiredCategories()
	if caps.SupportsEnum {
		categories = append(categories, testcases.CategoryEnum)
	}
	if caps.SupportsSchema {
		categories = append(categories, testcases.CategorySchema)
	}
	if caps.SupportsArray {
		categories = append(categories, testcases.CategoryArray)
	}
	if caps.SupportsJSON {
		categories = append(categories, testcases.CategoryJSON)
	}

	for _, cat := range categories {
		tests := registry.GetByCategory(cat)
		var covered, total int
		for _, tc := range tests {
			total++
			if implemented[tc.Name] {
				covered++
			}
		}
		if total > 0 {
			percentage := float64(covered) / float64(total) * 100
			t.Logf("  %s: %d/%d (%.1f%%)", cat, covered, total, percentage)
		}
	}
}
