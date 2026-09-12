package opts_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
)

func TestValidateOptsEmitIterators(t *testing.T) {
	limit := int32(1)
	base := func() *opts.Options {
		return &opts.Options{
			Package:             "db",
			QueryParameterLimit: &limit,
		}
	}

	if err := opts.ValidateOpts(base()); err != nil {
		t.Fatalf("valid defaults: %v", err)
	}

	cases := []struct {
		name    string
		mutate  func(*opts.Options)
		wantErr string
	}{
		{
			name: "invalid scope",
			mutate: func(o *opts.Options) {
				o.EmitIterators = true
				o.IteratorScope = "nope"
				o.IteratorStyle = "seq2"
				o.IteratorStart = "lazy"
			},
			wantErr: "iterator_scope",
		},
		{
			name: "invalid style",
			mutate: func(o *opts.Options) {
				o.EmitIterators = true
				o.IteratorScope = "global"
				o.IteratorStyle = "chan"
				o.IteratorStart = "lazy"
			},
			wantErr: "iterator_style",
		},
		{
			name: "invalid start",
			mutate: func(o *opts.Options) {
				o.EmitIterators = true
				o.IteratorScope = "global"
				o.IteratorStyle = "seq2"
				o.IteratorStart = "now"
			},
			wantErr: "iterator_start",
		},
		{
			name: "valid explicit_only",
			mutate: func(o *opts.Options) {
				o.EmitIterators = true
				o.IteratorScope = "explicit_only"
				o.IteratorStyle = "seq2"
				o.IteratorStart = "lazy"
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := base()
			tc.mutate(o)
			err := opts.ValidateOpts(o)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}
