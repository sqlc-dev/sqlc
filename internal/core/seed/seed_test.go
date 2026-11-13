package seed

import (
	"strings"
	"testing"
	"testing/fstest"
)

// Dialect seeds are embedded, so the CLI cannot supply invalid policy metadata.
func TestTypeCombinationStrategy(t *testing.T) {
	tests := []struct {
		name     string
		settings string
		wantErr  string
	}{
		{"disabled", `{"dialect":"postgresql"}`, ""},
		{"mssql", `{"dialect":"mssql","set_operations":{"strategy":"mssql"}}`, ""},
		{"missing", `{"dialect":"mssql","set_operations":{}}`, `unsupported type-combination strategy ""`},
		{"unknown", `{"dialect":"mssql","set_operations":{"strategy":"unknown"}}`, `unsupported type-combination strategy "unknown"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadSettings(fstest.MapFS{
				SettingsFile: &fstest.MapFile{Data: []byte(tt.settings)},
			})
			if tt.wantErr == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}
