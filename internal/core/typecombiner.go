package core

// TypeCombinationContext distinguishes constructs whose coercion rules may
// differ even within one dialect.
type TypeCombinationContext string

const (
	CombineUnion     TypeCombinationContext = "union"
	CombineIntersect TypeCombinationContext = "intersect"
	CombineExcept    TypeCombinationContext = "except"
	CombineCase      TypeCombinationContext = "case"
	CombineCoalesce  TypeCombinationContext = "coalesce"
	CombineGreatest  TypeCombinationContext = "greatest"
	CombineLeast     TypeCombinationContext = "least"
	CombineArray     TypeCombinationContext = "array"
)

type TypeCombinationRequest struct {
	Context TypeCombinationContext
	// FunctionName preserves a lowered function's identity as qualified name
	// components, not a result-type instruction. Nil means the ordinary
	// construct identified by Context.
	FunctionName []string
	// Candidates are in expression order, with nil for an untyped value.
	// Top-level result nullability is handled by the analyzer.
	Candidates []*TypeExpr
}

type TypeCombinationResult struct {
	// Handled distinguishes an untyped result from an unsupported context.
	// When false, the analyzer retains its existing inference behavior.
	Handled bool
	Type    *TypeExpr
}

// TypeCombiner is supplied by a dialect, not selected by shared analysis.
// Implementations must support concurrent calls and not mutate their inputs.
type TypeCombiner interface {
	Combine(TypeCombinationRequest) (TypeCombinationResult, error)
}
