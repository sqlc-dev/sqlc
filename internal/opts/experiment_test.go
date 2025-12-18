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
			name:  "enable analyzerv2",
			input: "analyzerv2",
			want:  Experiment{AnalyzerV2: true},
		},
		{
			name:  "disable analyzerv2",
			input: "noanalyzerv2",
			want:  Experiment{AnalyzerV2: false},
		},
		{
			name:  "enable then disable analyzerv2",
			input: "analyzerv2,noanalyzerv2",
			want:  Experiment{AnalyzerV2: false},
		},
		{
			name:  "analyzerv2 case insensitive",
			input: "AnalyzerV2",
			want:  Experiment{AnalyzerV2: true},
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
			name: "analyzerv2 enabled",
			exp:  Experiment{AnalyzerV2: true},
			want: []string{"analyzerv2"},
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
		{
			name: "analyzerv2 enabled",
			exp:  Experiment{AnalyzerV2: true},
			want: "analyzerv2",
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
		{
			name:  "analyzerv2 lowercase",
			input: "analyzerv2",
			want:  true,
		},
		{
			name:  "analyzerv2 mixed case",
			input: "AnalyzerV2",
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
