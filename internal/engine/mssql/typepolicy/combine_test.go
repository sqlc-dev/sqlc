package typepolicy

import (
	"strings"
	"testing"
)

// Embedded dialect metadata cannot be replaced through the CLI.
func TestDecodeRules(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{"missing", "", "EOF"},
		{"empty", "{}", "precedence must not be empty"},
		{"null", "null", "precedence must not be empty"},
		{"unknown field", `{"strategy":"mssql","precedence":["int"]}`, `unknown field "strategy"`},
		{"invalid shape", `{"precedence":"int"}`, "cannot unmarshal"},
		{"valid", `{"precedence":["int"]}`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeRules([]byte(tt.data))
			if tt.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v, want error containing %q", err, tt.want)
			}
		})
	}
}
