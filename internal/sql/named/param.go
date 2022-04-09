package named

// Nullability represents the nullability of a named parameter.
// The nullability can be:
// 1. unspecified
// 2. inferred
// 3. user-defined
// A user-specified nullability carries a higher precedence than an inferred one
//
// The representation is such that you can bitwise OR together Nullability types to
// combine them together.
type Nullability int

const (
	NullUnspecified Nullability = 0b0000
	InferredNull    Nullability = 0b0001
	InferredNotNull Nullability = 0b0010
	Nullable        Nullability = 0b0100
	NotNullable     Nullability = 0b1000
)

// String implements the Stringer interface
func (n Nullability) String() string {
	switch n {
	case NullUnspecified:
		return "NullUnspecified"
	case InferredNull:
		return "InferredNull"
	case InferredNotNull:
		return "InferredNotNull"
	case Nullable:
		return "Nullable"
	case NotNullable:
		return "NotNullable"
	default:
		return "NullInvalid"
	}
}

// Param represents a input argument to the query which can be specified using:
// - positional parameters           $1
// - named parameter operator        @param
// - named parameter function calls  sqlc.arg(param)
type Param struct {
	name        string
	nullability Nullability
}

// NewUnspecifiedParam builds a new params with unspecified nullability
func NewUnspecifiedParam(name string) Param {
	return Param{name: name, nullability: NullUnspecified}
}

// NewInferredParam builds a new params with inferred nullability
func NewInferredParam(name string, notNull bool) Param {
	if notNull {
		return Param{name: name, nullability: InferredNotNull}
	}

	return Param{name: name, nullability: InferredNull}
}

// NewUserDefinedParam creates a new param with the user specified
// by the end user
func NewUserDefinedParam(name string, notNull bool) Param {
	if notNull {
		return Param{name: name, nullability: NotNullable}
	}

	return Param{name: name, nullability: Nullable}
}

// Name is the user defined name to use for this parameter
func (p Param) Name() string {
	return p.name
}

// is checks if this params object has the specified nullability bit set
func (p Param) is(n Nullability) bool {
	return (p.nullability & n) == n
}

// NonNull determines whether this param should be "not null" in its current state
func (p Param) NotNull() bool {
	const nullable = false
	const notNull = true

	if p.is(NotNullable) {
		return notNull
	}

	if p.is(Nullable) {
		return nullable
	}

	if p.is(InferredNotNull) {
		return notNull
	}

	if p.is(InferredNull) {
		return nullable
	}

	// This param is unspecified, so by default we choose nullable
	// which matches the default behavior of most databases
	return nullable
}

// Combine creates a new param from 2 partially specified params
// If the parameters have different names, the first is preferred
func Combine(a, b Param) Param {
	name := a.name
	if name == "" {
		name = b.name
	}

	return Param{name: name, nullability: a.nullability | b.nullability}
}
