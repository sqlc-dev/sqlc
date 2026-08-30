package opts

import (
	"strings"
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

func TestMatchesColumnTimestamptzAliases(t *testing.T) {
	t.Parallel()

	parseOverride := func(t *testing.T, dbType string, nullable bool) Override {
		t.Helper()
		o := Override{
			DBType:   dbType,
			Nullable: nullable,
			GoType:   GoType{Spec: "*time.Time"},
		}
		if err := o.parse(nil); err != nil {
			t.Fatalf("override parsing failed: %s", err)
		}
		return o
	}

	column := func(typeName string, nullable bool) *plugin.Column {
		typ := &plugin.Identifier{Name: typeName}
		if schema, name, ok := strings.Cut(typeName, "."); ok && schema == "pg_catalog" {
			typ = &plugin.Identifier{Schema: schema, Name: name}
		}
		return &plugin.Column{
			Type:    typ,
			NotNull: !nullable,
		}
	}

	for _, test := range []struct {
		name       string
		override   Override
		column     *plugin.Column
		wantMatch  bool
	}{
		{
			name:      "timestamptz override matches timestamptz column",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("timestamptz", true),
			wantMatch: true,
		},
		{
			name:      "timestamptz override matches timestamp with time zone column",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("timestamp with time zone", true),
			wantMatch: true,
		},
		{
			name:      "timestamptz override matches pg_catalog.timestamptz column",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("pg_catalog.timestamptz", true),
			wantMatch: true,
		},
		{
			name:      "pg_catalog.timestamptz override matches timestamptz column",
			override:  parseOverride(t, "pg_catalog.timestamptz", true),
			column:    column("timestamptz", true),
			wantMatch: true,
		},
		{
			name:      "pg_catalog.timestamptz override matches timestamp with time zone column",
			override:  parseOverride(t, "pg_catalog.timestamptz", true),
			column:    column("timestamp with time zone", true),
			wantMatch: true,
		},
		{
			name:      "timestamp with time zone override matches timestamptz column",
			override:  parseOverride(t, "timestamp with time zone", true),
			column:    column("timestamptz", true),
			wantMatch: true,
		},
		{
			name:      "timestamptz override does not match not-null column",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("timestamptz", false),
			wantMatch: false,
		},
		{
			name:      "timestamptz override does not match timestamp without time zone",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("timestamp", true),
			wantMatch: false,
		},
		{
			name:      "timestamptz override does not match timestamp without time zone long form",
			override:  parseOverride(t, "timestamptz", true),
			column:    column("timestamp without time zone", true),
			wantMatch: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.override.MatchesColumn(test.column)
			if got != test.wantMatch {
				t.Errorf("MatchesColumn() = %v, want %v", got, test.wantMatch)
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
