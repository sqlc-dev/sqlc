package config

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
  "version": "1",
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
			`yaml: unmarshal errors:
  line 3: field foo not found in type config.V1GenerateSettings`,
			unknownFields,
		},
	} {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfig(strings.NewReader(tt.json))
			if err == nil {
				t.Fatalf("expected err; got nil")
			}
			if diff := cmp.Diff(err.Error(), tt.err); diff != "" {
				t.Errorf("differed (-want +got):\n%s", diff)
			}
		})
	}
}

const validConfigOne = `{
  "version": "1"
  "packages": []
}`

func FuzzConfig(f *testing.F) {
	f.Add(validConfigOne)
	f.Fuzz(func(t *testing.T, orig string) {
		ParseConfig(strings.NewReader(orig))
	})
}

func TestInvalidConfig(t *testing.T) {
	err := Validate(&Config{
		SQL: []SQL{{
			Gen: SQLGen{
				Go: &SQLGo{
					EmitMethodsWithDBArgument: true,
					EmitPreparedQueries:       true,
				},
			},
		}}})
	if err == nil {
		t.Errorf("expected err; got nil")
	}
}

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
	} {
		tt := test
		t.Run(tt.override.GoType.Spec, func(t *testing.T) {
			if err := tt.override.Parse(); err != nil {
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
			err := tt.override.Parse()
			if err == nil {
				t.Fatalf("expected pars to fail; got nil")
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
		o.Parse()
	})
}
