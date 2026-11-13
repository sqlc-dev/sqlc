package analyzer

import (
	"errors"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/core"
)

type combinerFunc func(core.TypeCombinationRequest) (core.TypeCombinationResult, error)

func (f combinerFunc) Combine(req core.TypeCombinationRequest) (core.TypeCombinationResult, error) {
	return f(req)
}

// A deliberately declining/failing engine cannot be selected through the CLI.
func TestTypeCombinerContract(t *testing.T) {
	failure := errors.New("policy failure")
	for _, name := range []string{"absent", "declined", "untyped result", "error"} {
		t.Run(name, func(t *testing.T) {
			types := []exprType{{}, {expr: core.ParseTypeExpr("integer")}, {expr: core.ParseTypeExpr("text")}}
			a := &analyzer{}
			if name != "absent" {
				a.combiner = combinerFunc(func(req core.TypeCombinationRequest) (core.TypeCombinationResult, error) {
					if req.Context != core.CombineCoalesce || len(req.FunctionName) != 2 ||
						req.FunctionName[0] != "schema" || req.FunctionName[1] != "fn" ||
						len(req.Candidates) != 3 ||
						req.Candidates[0] != nil || req.Candidates[1].Name != "integer" ||
						req.Candidates[2].Name != "text" {
						t.Fatalf("candidate context/order not preserved: %+v", req)
					}
					switch name {
					case "error":
						return core.TypeCombinationResult{}, failure
					case "untyped result":
						return core.TypeCombinationResult{Handled: true}, nil
					default:
						return core.TypeCombinationResult{}, nil
					}
				})
			}
			got, err := a.commonType(core.TypeCombinationRequest{
				Context: core.CombineCoalesce, FunctionName: []string{"schema", "fn"},
			}, types)
			if name == "error" {
				if !errors.Is(err, failure) {
					t.Fatalf("got %v, want policy failure", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if name == "untyped result" {
				if got.typeOID != 0 || got.expr != nil {
					t.Fatalf("handled untyped result fell back to %+v", got)
				}
			} else if got.expr == nil || got.expr.Name != "integer" {
				t.Fatalf("first-typed fallback changed: %+v", got)
			}
		})
	}
}
