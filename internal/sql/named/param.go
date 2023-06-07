package named

// nullability represents the nullability of a named parameter.
// The nullability can be:
// 1. unspecified
// 2. inferred
// 3. user-defined
// A user-specified nullability carries a higher precedence than an inferred one
//
// The representation is such that you can bitwise OR together nullability types to
// combine them together.
type nullability int

const (
	nullUnspecified nullability = 0b0000
	inferredNull    nullability = 0b0001
	inferredNotNull nullability = 0b0010
	nullable        nullability = 0b0100
	notNullable     nullability = 0b1000
)

// String implements the Stringer interface
func (n nullability) String() string {
	switch n {
	case nullUnspecified:
		return "NullUnspecified"
	case inferredNull:
		return "InferredNull"
	case inferredNotNull:
		return "InferredNotNull"
	case nullable:
		return "Nullable"
	case notNullable:
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
	nullability nullability
	isSqlcSlice bool
}

// NewParam builds a new params with unspecified nullability
func NewParam(name string) Param {
	return Param{name: name, nullability: nullUnspecified}
}

// NewInferredParam builds a new params with inferred nullability
func NewInferredParam(name string, notNull bool) Param {
	if notNull {
		return Param{name: name, nullability: inferredNotNull}
	}

	return Param{name: name, nullability: inferredNull}
}

// NewUserNullableParam is a parameter that has been overridden
// by the user to be nullable.
func NewUserNullableParam(name string) Param {
	return Param{name: name, nullability: nullable}
}

// NewSqlcSlice is a sqlc.slice() parameter.
func NewSqlcSlice(name string) Param {
	return Param{name: name, nullability: nullUnspecified, isSqlcSlice: true}
}

// Name is the user defined name to use for this parameter
func (p Param) Name() string {
	return p.name
}

// is checks if this params object has the specified nullability bit set
func (p Param) is(n nullability) bool {
	return (p.nullability & n) == n
}

// NonNull determines whether this param should be "not null" in its current state
func (p Param) NotNull() bool {
	const null = false
	const notNull = true

	if p.is(notNullable) {
		return notNull
	}

	if p.is(nullable) {
		return null
	}

	if p.is(inferredNotNull) {
		return notNull
	}

	if p.is(inferredNull) {
		return null
	}

	// This param is unspecified, so by default we choose nullable
	// which matches the default behavior of most databases
	return null
}

// IsSlice returns whether this param is a sqlc.slice() param.
func (p Param) IsSqlcSlice() bool {
	return p.isSqlcSlice
}

// mergeParam creates a new param from 2 partially specified params
// If the parameters have different names, the first is preferred
func mergeParam(a, b Param) Param {
	name := a.name
	if name == "" {
		name = b.name
	}

	return Param{
		name:        name,
		nullability: a.nullability | b.nullability,
		isSqlcSlice: a.isSqlcSlice || b.isSqlcSlice,
	}
}
