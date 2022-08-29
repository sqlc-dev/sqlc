package named

import "testing"

func TestParamSet_Add(t *testing.T) {
	t.Parallel()

	type test struct {
		pset     *ParamSet
		param    Param
		expected int
	}

	named := NewParamSet(nil, true)
	populatedNamed := NewParamSet(map[int]bool{1: true, 2: true, 4: true, 5: true, 6: true}, true)
	populatedUnnamed := NewParamSet(map[int]bool{1: true, 2: true, 4: true, 5: true, 6: true}, false)
	unnamed := NewParamSet(nil, false)
	p1 := NewParam("hello")
	p2 := NewParam("world")

	tests := []test{
		// First parameter should be 1
		{named, p1, 1},
		// Duplicate first parameters should be 1
		{named, p1, 1},
		// A new parameter receives a new parameter number
		{named, p2, 2},
		// An additional new parameter does _not_ receive a new
		{named, p2, 2},

		// First parameter should be 1
		{unnamed, p1, 1},
		// Duplicate first parameters should increment argn
		{unnamed, p1, 2},
		// A new parameter receives a new parameter number
		{unnamed, p2, 3},
		// An additional new parameter still does receive a new argn
		{unnamed, p2, 4},

		// First parameter of a pre-populated should be 3
		{populatedNamed, p1, 3},
		{populatedNamed, p1, 3},
		{populatedNamed, p2, 7},
		{populatedNamed, p2, 7},

		{populatedUnnamed, p1, 3},
		{populatedUnnamed, p1, 7},
		{populatedUnnamed, p2, 8},
		{populatedUnnamed, p2, 9},
	}

	for _, spec := range tests {
		actual := spec.pset.Add(spec.param)
		if actual != spec.expected {
			t.Errorf("ParamSet.Add(%s) expected %v; got %v", spec.param.name, spec.expected, actual)
		}
	}
}
