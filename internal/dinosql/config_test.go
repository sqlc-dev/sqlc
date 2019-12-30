package dinosql

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const missingVersion = `{
}`

const missingPackages = `{
  "version": "1"
}`

const unknownVersion = `{
  "version": "foo"
}`

const unknownFields = `{
  "foo": "bar"
}`

func TestBadConfigs(t *testing.T) {
	for _, test := range []struct {
		name string
		err  string
		json string
	}{
		{
			"missing version",
			"no version number",
			missingVersion,
		},
		{
			"missing packages",
			"no packages",
			missingPackages,
		},
		{
			"unknown version",
			"invalid version number",
			unknownVersion,
		},
		{
			"unknown fields",
			"json: unknown field \"foo\"",
			unknownFields,
		},
	} {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfigFile(strings.NewReader(tt.json))
			if err == nil {
				t.Fatalf("expected err; got nil")
			}
			if diff := cmp.Diff(err.Error(), tt.err); diff != "" {
				t.Errorf("differed (-want +got):\n%s", diff)
			}
		})
	}
}
