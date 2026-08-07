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
		{
			name:  "enable coreanalyzer",
			input: "coreanalyzer",
			want:  Experiment{CoreAnalyzer: true},
		},
		{
			name:  "disable coreanalyzer",
			input: "nocoreanalyzer",
			want:  Experiment{CoreAnalyzer: false},
		},
		{
			name:  "enable then disable coreanalyzer",
			input: "coreanalyzer,nocoreanalyzer",
			want:  Experiment{CoreAnalyzer: false},
		},
		{
			name:  "coreanalyzer case insensitive",
			input: "CoreAnalyzer",
			want:  Experiment{CoreAnalyzer: true},
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
		{
			name: "coreanalyzer enabled",
			exp:  Experiment{CoreAnalyzer: true},
			want: []string{"coreanalyzer"},
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
		{
			name:  "coreanalyzer lowercase",
			input: "coreanalyzer",
			want:  true,
		},
		{
			name:  "coreanalyzer mixed case",
			input: "CoreAnalyzer",
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
