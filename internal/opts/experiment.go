package opts

import (
	"os"
	"strings"
)

// The SQLCEXPERIMENT variable controls experimental features within sqlc. It
// is a comma-separated list of experiment names. Experiment names can be
// prefixed with "no" to explicitly disable them.
//
// This is modeled after Go's GOEXPERIMENT environment variable. For more
// information, see https://pkg.go.dev/internal/goexperiment
//
// Available experiments:
//
//	coreanalyzer - routes generate through the core catalog and analyzer
//
// Example usage:
//
//	SQLCEXPERIMENT=foo,bar      # enable foo and bar experiments
//	SQLCEXPERIMENT=nofoo        # explicitly disable foo experiment
//	SQLCEXPERIMENT=foo,nobar    # enable foo, disable bar

// Experiment holds the state of all experimental features.
// Add new experiments as boolean fields to this struct.
type Experiment struct {
	// CoreAnalyzer routes generate through the core catalog and analyzer
	// instead of each engine's own analysis path.
	CoreAnalyzer bool
}

// ExperimentFromEnv returns an Experiment initialized from the SQLCEXPERIMENT
// environment variable.
func ExperimentFromEnv() Experiment {
	return ExperimentFromString(os.Getenv("SQLCEXPERIMENT"))
}

// ExperimentFromString parses a comma-separated list of experiment names
// and returns an Experiment with the appropriate flags set.
//
// Experiment names can be prefixed with "no" to explicitly disable them.
// Unknown experiment names are silently ignored.
func ExperimentFromString(val string) Experiment {
	e := Experiment{}
	if val == "" {
		return e
	}

	for name := range strings.SplitSeq(val, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		// Check if this is a negation (noFoo)
		enabled := true
		if strings.HasPrefix(strings.ToLower(name), "no") && len(name) > 2 {
			// Could be a negation, check if the rest is a valid experiment
			possibleExp := name[2:]
			if isKnownExperiment(possibleExp) {
				name = possibleExp
				enabled = false
			}
			// If not a known experiment, treat "no..." as a potential experiment name itself
		}

		setExperiment(&e, name, enabled)
	}

	return e
}

// isKnownExperiment returns true if the given name (case-insensitive) is a
// known experiment.
func isKnownExperiment(name string) bool {
	switch strings.ToLower(name) {
	case "coreanalyzer":
		return true
	default:
		return false
	}
}

// setExperiment sets the experiment flag with the given name to the given value.
func setExperiment(e *Experiment, name string, enabled bool) {
	switch strings.ToLower(name) {
	case "coreanalyzer":
		e.CoreAnalyzer = enabled
	}
}

// Enabled returns a slice of all enabled experiment names.
func (e Experiment) Enabled() []string {
	var enabled []string
	if e.CoreAnalyzer {
		enabled = append(enabled, "coreanalyzer")
	}
	return enabled
}

// String returns a comma-separated list of enabled experiments.
func (e Experiment) String() string {
	return strings.Join(e.Enabled(), ",")
}
