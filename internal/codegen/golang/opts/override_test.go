package opts

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
