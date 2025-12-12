package opts

import "testing"

func TestExperimentFromString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Experiment
	}{
		{
			name:  "empty string",
			input: "",
			want:  Experiment{},
		},
		{
			name:  "whitespace only",
			input: "   ",
			want:  Experiment{},
		},
		{
			name:  "unknown experiment",
			input: "unknownexperiment",
			want:  Experiment{},
		},
		{
			name:  "multiple unknown experiments",
			input: "foo,bar,baz",
			want:  Experiment{},
		},
		{
			name:  "unknown with no prefix",
			input: "nounknown",
			want:  Experiment{},
		},
		{
			name:  "whitespace around experiments",
			input: " foo , bar , baz ",
			want:  Experiment{},
		},
		{
			name:  "empty items in list",
			input: "foo,,bar",
			want:  Experiment{},
		},
		// Add tests for specific experiments as they are introduced.
		// Example:
		// {
		// 	name:  "enable newparser",
		// 	input: "newparser",
		// 	want:  Experiment{NewParser: true},
		// },
		// {
		// 	name:  "disable newparser",
		// 	input: "nonewparser",
		// 	want:  Experiment{NewParser: false},
		// },
		// {
		// 	name:  "enable then disable",
		// 	input: "newparser,nonewparser",
		// 	want:  Experiment{NewParser: false},
		// },
		// {
		// 	name:  "case insensitive",
		// 	input: "NewParser,NONEWPARSER",
		// 	want:  Experiment{NewParser: false},
		// },
		{
			name:  "enable pglite",
			input: "pglite",
			want:  Experiment{PGLite: true},
		},
		{
			name:  "disable pglite",
			input: "nopglite",
			want:  Experiment{PGLite: false},
		},
		{
			name:  "pglite enable then disable",
			input: "pglite,nopglite",
			want:  Experiment{PGLite: false},
		},
		{
			name:  "pglite case insensitive",
			input: "PGLite",
			want:  Experiment{PGLite: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExperimentFromString(tt.input)
			if got != tt.want {
				t.Errorf("ExperimentFromString(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExperimentEnabled(t *testing.T) {
	tests := []struct {
		name string
		exp  Experiment
		want []string
	}{
		{
			name: "no experiments enabled",
			exp:  Experiment{},
			want: nil,
		},
		// Add tests for specific experiments as they are introduced.
		// Example:
		// {
		// 	name: "newparser enabled",
		// 	exp:  Experiment{NewParser: true},
		// 	want: []string{"newparser"},
		// },
		{
			name: "pglite enabled",
			exp:  Experiment{PGLite: true},
			want: []string{"pglite"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.exp.Enabled()
			if len(got) != len(tt.want) {
				t.Errorf("Experiment.Enabled() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Experiment.Enabled()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExperimentString(t *testing.T) {
	tests := []struct {
		name string
		exp  Experiment
		want string
	}{
		{
			name: "no experiments",
			exp:  Experiment{},
			want: "",
		},
		// Add tests for specific experiments as they are introduced.
		// Example:
		// {
		// 	name: "newparser enabled",
		// 	exp:  Experiment{NewParser: true},
		// 	want: "newparser",
		// },
		{
			name: "pglite enabled",
			exp:  Experiment{PGLite: true},
			want: "pglite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.exp.String()
			if got != tt.want {
				t.Errorf("Experiment.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsKnownExperiment(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "unknown experiment",
			input: "unknown",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		// Add tests for specific experiments as they are introduced.
		// Example:
		// {
		// 	name:  "newparser lowercase",
		// 	input: "newparser",
		// 	want:  true,
		// },
		// {
		// 	name:  "newparser mixed case",
		// 	input: "NewParser",
		// 	want:  true,
		// },
		{
			name:  "pglite lowercase",
			input: "pglite",
			want:  true,
		},
		{
			name:  "pglite mixed case",
			input: "PGLite",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isKnownExperiment(tt.input)
			if got != tt.want {
				t.Errorf("isKnownExperiment(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
