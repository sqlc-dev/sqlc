package golang

import (
	"reflect"
	"testing"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func TestPutOutColumns_ForZeroColumns(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{
			cmd:  metadata.CmdExec,
			want: false,
		},
		{
			cmd:  metadata.CmdExecResult,
			want: false,
		},
		{
			cmd:  metadata.CmdExecRows,
			want: false,
		},
		{
			cmd:  metadata.CmdExecLastId,
			want: false,
		},
		{
			cmd:  metadata.CmdMany,
			want: true,
		},
		{
			cmd:  metadata.CmdOne,
			want: true,
		},
		{
			cmd:  metadata.CmdCopyFrom,
			want: false,
		},
		{
			cmd:  metadata.CmdBatchExec,
			want: false,
		},
		{
			cmd:  metadata.CmdBatchMany,
			want: true,
		},
		{
			cmd:  metadata.CmdBatchOne,
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.cmd, func(t *testing.T) {
			query := &plugin.Query{
				Cmd:     tc.cmd,
				Columns: []*plugin.Column{},
			}
			got := putOutColumns(query)
			if got != tc.want {
				t.Errorf("putOutColumns failed. want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestPutOutColumns_AlwaysTrueWhenQueryHasColumns(t *testing.T) {
	query := &plugin.Query{
		Cmd:     metadata.CmdMany,
		Columns: []*plugin.Column{{}},
	}
	if putOutColumns(query) != true {
		t.Error("should be true when we have columns")
	}
}

func Test_removeUnused(t *testing.T) {
	type args struct {
		enums   []Enum
		structs []Struct
		queries []Query
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]Enum
		want1 map[string]Struct
	}{
		{
			name: "remove unused structs and enums",
			args: args{
				enums: []Enum{
					{
						Name: "enum1",
					},
					{
						Name: "enum2",
					},
				},
				structs: []Struct{
					{
						Name: "struct1",
					},
					{
						Name: "struct2",
					},
					{
						Name: "struct3",
					},
					{
						Name: "struct4",
						Fields: []Field{
							{
								Type: "enum1",
							},
							{
								Type: "notenum",
							},
						},
					},
					{
						Name: "struct5",
					},
				},
				queries: []Query{
					{
						Cmd: metadata.CmdOne,
						Ret: QueryValue{
							Struct: &Struct{
								Name: "struct1",
							},
						},
						Arg: QueryValue{
							Struct: &Struct{
								Name: "struct2",
							},
						},
					},
					{
						Cmd: metadata.CmdOne,
						Ret: QueryValue{
							Struct: &Struct{
								Name: "struct3",
							},
						},
						Arg: QueryValue{
							Struct: &Struct{
								Name: "struct4",
								Fields: []Field{
									{
										Type: "enum1",
									},
									{
										Type: "notenum",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]Enum{
				"enum1": {Name: "enum1"},
			},
			want1: map[string]Struct{
				"struct1": {
					Name: "struct1",
				},
				"struct2": {
					Name: "struct2",
				},
				"struct3": {
					Name: "struct3",
				},
				"struct4": {
					Name: "struct4",
					Fields: []Field{
						{
							Type: "enum1",
						},
						{
							Type: "notenum",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := removeUnused(tt.args.enums, tt.args.structs, tt.args.queries)
			gotMap := map[string]Enum{}
			got1Map := map[string]Struct{}
			for _, enum := range got {
				gotMap[enum.Name] = enum
			}
			for _, s := range got1 {
				got1Map[s.Name] = s
			}
			if !reflect.DeepEqual(gotMap, tt.want) {
				t.Errorf("removeUnused() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1Map, tt.want1) {
				t.Errorf("removeUnused() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
