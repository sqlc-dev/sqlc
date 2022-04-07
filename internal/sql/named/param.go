package named

// Nullability represents the nullability of a named parameter The
// representation is such that you can bitwise OR together Nullability types to
// combine them
// For example:
// - NullUnspecified | Nullable = Nullable
// - NonNullable     | Nullable = NullInvalid
type Nullability int

const (
	NullUnspecified Nullability = 0b00
	Nullable        Nullability = 0b01
	NotNullable     Nullability = 0b10
	NullInvalid     Nullability = 0b11
)

// String implements the Stringer interface
func (n Nullability) String() string {
	switch n {
	case NullUnspecified:
		return "NullUnspecified"
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

// NewParam creates a new named param with the given nullability
func NewParam(name string, notNull bool) Param {
	if notNull {
		return Param{name: name, nullability: NotNullable}
	}

	return Param{name: name, nullability: Nullable}
}

// Name is the user defined name to use for this parameter
func (p Param) Name() string {
	return p.name
}

// Nullability retrieves the nullability status of this param
func (p Param) Nullability() Nullability {
	return p.nullability
}

// NonNull determines whether this param is NonNull
func (p Param) NotNull() bool {
	return (p.nullability & NotNullable) > 0
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
