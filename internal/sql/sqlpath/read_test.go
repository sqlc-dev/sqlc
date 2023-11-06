package sqlpath

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// Returns a list of SQL files from given paths.
func TestReturnsListOfSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file1.sql", "testdata/file2.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{"testdata/file1.sql", "testdata/file2.sql"}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v, %v", expected, result, cmp.Diff(expected, result))
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestReturnsNilListWhenNoSQLFilesFound(t *testing.T) {
	// Arrange
	paths := []string{"testdata/extra.txt"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v, %v", expected, result, cmp.Diff(expected, result))
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestIgnoresHiddenFilesWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/.hidden.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestIgnoresNonSQLFilesWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/extra.txt"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestExcludesSQLFilesEndingWithDownSQLWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file1.sql", "testdata/file3.down.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{"testdata/file1.sql"}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestReturnsErrorWhenPathDoesNotExist(t *testing.T) {
	// Arrange
	paths := []string{"non_existent_path"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else {
		expectedError := fmt.Errorf("path non_existent_path does not exist")
		if !cmp.Equal(err.Error(), expectedError.Error()) {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}
	}
}

func TestReturnsErrorWhenDirectoryCannotBeRead(t *testing.T) {
	// Arrange
	paths := []string{"testdata/unreadable"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else {
		expectedError := fmt.Errorf("open testdata/unreadable: permission denied")
		if !cmp.Equal(err.Error(), expectedError.Error()) {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}
	}
}

func TestDoesNotIncludesSQLFilesWithUppercaseExtension(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file4.SQL"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestIncludesSQLFilesWithLeadingDotsInDirectoryName(t *testing.T) {
	// Arrange
	paths := []string{"./testdata/.hiddendir/file1.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{"./testdata/.hiddendir/file1.sql"}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestPathIsSymlink(t *testing.T) {
	// Arrange
	paths := []string{"testdata/symlink", "testdata/file1.symlink.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{
		"testdata/symlink/file1.sql",
		"testdata/symlink/file1.symlink.sql",
		"testdata/symlink/file2.sql",
		"testdata/file1.symlink.sql",
	}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}
