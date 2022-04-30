package named

import "testing"

func TestMergeParamNullability(t *testing.T) {
	type test struct {
		a       Param
		b       Param
		notNull bool
		message string
	}

	name := "name"
	unspec := NewParam(name)
	inferredNotNull := NewInferredParam(name, true)
	inferredNull := NewInferredParam(name, false)
	userDefNull := NewUserNullableParam(name)

	const notNull = true
	const null = false

	tests := []test{
		// Unspecified nullability parameter works
		{unspec, inferredNotNull, notNull, "Unspec + inferred(not null) = not null"},
		{unspec, inferredNull, null, "Unspec + inferred(not null) = null"},
		{unspec, userDefNull, null, "Unspec + userdef(null) = null"},

		// Inferred nullability agreeing with user defined nullabilty
		{inferredNull, userDefNull, null, "inferred(null) + userdef(null) = null"},

		// Inferred nullability disagreeing with user defined nullabilty
		{inferredNotNull, userDefNull, null, "inferred(not null) + userdef(null) = null"},
	}

	for _, spec := range tests {
		a := spec.a
		b := spec.b
		actual := mergeParam(a, b).NotNull()
		expected := spec.notNull
		if actual != expected {
			t.Errorf("Combine(%s,%s) expected %v; got %v", a.nullability, b.nullability, expected, actual)
		}

		// We have already tried Combine(a, b) the same result should be true for Combine(b, a)
		actual = mergeParam(b, a).NotNull()
		if actual != expected {
			t.Errorf("Combine(%s,%s) expected %v; got %v", b.nullability, a.nullability, expected, actual)
		}
	}
}

func TestMergeParamName(t *testing.T) {
	type test struct {
		a    Param
		b    Param
		name string
	}

	a := NewParam("a")
	b := NewParam("b")
	blank := NewParam("")

	tests := []test{
		// should prefer the first param's name if both specified
		{a, b, "a"},
		{b, a, "b"},

		// should prefer non-blank names
		{a, blank, "a"},
		{blank, a, "a"},
	}

	for _, spec := range tests {
		a := spec.a
		b := spec.b
		actual := mergeParam(a, b).Name()
		expected := spec.name
		if actual != expected {
			t.Errorf("Combine(%s,%s) expected %v; got %v", a.name, b.name, expected, actual)
		}
	}
}
