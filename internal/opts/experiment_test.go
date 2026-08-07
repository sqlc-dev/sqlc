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
	exp := Experiment{}
	if got := exp.Enabled(); len(got) != 0 {
		t.Errorf("Experiment.Enabled() = %v, want none", got)
	}
}

func TestExperimentString(t *testing.T) {
	exp := Experiment{}
	if got := exp.String(); got != "" {
		t.Errorf("Experiment.String() = %q, want %q", got, "")
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
