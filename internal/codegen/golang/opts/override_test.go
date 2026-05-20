package opts

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func TestTypeOverrides(t *testing.T) {
	for _, test := range []struct {
		override Override
		pkg      string
		typeName string
		basic    bool
	}{
		{
			Override{
				DBType: "uuid",
				GoType: GoType{Spec: "github.com/segmentio/ksuid.KSUID"},
			},
			"github.com/segmentio/ksuid",
			"ksuid.KSUID",
			false,
		},
		// TODO: Add test for struct pointers
		//
		// {
		// 	Override{
		// 		DBType: "uuid",
		// 		GoType:       "github.com/segmentio/*ksuid.KSUID",
		// 	},
		// 	"github.com/segmentio/ksuid",
		// 	"*ksuid.KSUID",
		// 	false,
		// },
		{
			Override{
				DBType: "citext",
				GoType: GoType{Spec: "string"},
			},
			"",
			"string",
			true,
		},
		{
			Override{
				DBType: "timestamp",
				GoType: GoType{Spec: "time.Time"},
			},
			"time",
			"time.Time",
			false,
		},
	} {
		tt := test
		t.Run(tt.override.GoType.Spec, func(t *testing.T) {
			if err := tt.override.parse(nil); err != nil {
				t.Fatalf("override parsing failed; %s", err)
			}
			if diff := cmp.Diff(tt.pkg, tt.override.GoImportPath); diff != "" {
				t.Errorf("package mismatch;\n%s", diff)
			}
			if diff := cmp.Diff(tt.typeName, tt.override.GoTypeName); diff != "" {
				t.Errorf("type name mismatch;\n%s", diff)
			}
			if diff := cmp.Diff(tt.basic, tt.override.GoBasicType); diff != "" {
				t.Errorf("basic mismatch;\n%s", diff)
			}
		})
	}
	for _, test := range []struct {
		override Override
		err      string
	}{
		{
			Override{
				DBType: "uuid",
				GoType: GoType{Spec: "Pointer"},
			},
			"Package override `go_type` specifier \"Pointer\" is not a Go basic type e.g. 'string'",
		},
		{
			Override{
				DBType: "uuid",
				GoType: GoType{Spec: "untyped rune"},
			},
			"Package override `go_type` specifier \"untyped rune\" is not a Go basic type e.g. 'string'",
		},
	} {
		tt := test
		t.Run(tt.override.GoType.Spec, func(t *testing.T) {
			err := tt.override.parse(nil)
			if err == nil {
				t.Fatalf("expected parse to fail; got nil")
			}
			if diff := cmp.Diff(tt.err, err.Error()); diff != "" {
				t.Errorf("error mismatch;\n%s", diff)
			}
		})
	}
}

func FuzzOverride(f *testing.F) {
	for _, spec := range []string{
		"string",
		"github.com/gofrs/uuid.UUID",
		"github.com/segmentio/ksuid.KSUID",
	} {
		f.Add(spec)
	}
	f.Fuzz(func(t *testing.T, s string) {
		o := Override{
			GoType: GoType{Spec: s},
		}
		o.parse(nil)
	})
}

func TestOverride_MatchesColumn(t *testing.T) {
	t.Parallel()
	type testCase struct {
		specName string
		override Override
		Column   *plugin.Column
		engine   string
		expected bool
	}

	testCases := []*testCase{
		{
			specName: "matches with pg_catalog in schema and name",
			override: Override{
				DBType:   "json",
				Nullable: false,
			},
			Column: &plugin.Column{
				Name: "data",
				Type: &plugin.Identifier{
					Schema: "pg_catalog",
					Name:   "json",
				},
				NotNull: true,
				IsArray: false,
			},
			engine:   "postgresql",
			expected: true,
		},
		{
			specName: "matches only with name",
			override: Override{
				DBType:   "json",
				Nullable: false,
			},
			Column: &plugin.Column{
				Name: "data",
				Type: &plugin.Identifier{
					Name: "json",
				},
				NotNull: true,
				IsArray: false,
			},
			engine:   "postgresql",
			expected: true,
		},
		{
			specName: "matches with pg_catalog in name",
			override: Override{
				DBType:   "json",
				Nullable: false,
			},
			Column: &plugin.Column{
				Name: "data",
				Type: &plugin.Identifier{
					Name: "pg_catalog.json",
				},
				NotNull: true,
				IsArray: false,
			},
			engine:   "postgresql",
			expected: true,
		},
	}

	for _, test := range testCases {
		tt := *test
		t.Run(tt.specName, func(t *testing.T) {
			result := tt.override.MatchesColumn(tt.Column, tt.engine)
			if result != tt.expected {
				t.Errorf("mismatch; got %v; want %v", result, tt.expected)
			}
			if tt.engine == "postgresql" && tt.expected == true {
				tt.override.DBType = "pg_catalog." + tt.override.DBType
				result = tt.override.MatchesColumn(test.Column, tt.engine)
				if !result {
					t.Errorf("mismatch; got %v; want %v", result, tt.expected)
				}
			}

		})

	}
}
