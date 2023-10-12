package named

// ParamSet represents a set of parameters for a single query
type ParamSet struct {
	// does this engine support named parameters?
	hasNamedSupport bool
	// the set of currently tracked named parameters
	namedParams map[string]Param
	// the locations of each of the named parameters
	namedLocs map[string][]int
	// a map of positions currently used
	positionToName map[int]string
	// argn keeps track of the last checked positional parameter used
	argn int
}

// Return the name for a given parameter number and a boolean indicating if it
// was found.
func (p *ParamSet) NameFor(idx int) (string, bool) {
	name, ok := p.positionToName[idx]
	return name, ok
}

func (p *ParamSet) nextArgNum() int {
	for {
		if _, ok := p.positionToName[p.argn]; !ok {
			return p.argn
		}

		p.argn++
	}
}

// Add adds a parameter to this set and returns the numbered location used for it
func (p *ParamSet) Add(param Param) int {
	name := param.name
	existing, ok := p.namedParams[name]

	p.namedParams[name] = mergeParam(existing, param)
	if ok && p.hasNamedSupport {
		return p.namedLocs[name][0]
	}

	argn := p.nextArgNum()
	p.positionToName[argn] = name
	p.namedLocs[name] = append(p.namedLocs[name], argn)
	return argn
}

// FetchMerge fetches an indexed parameter, and merges `mergeP` into it
// Returns: the merged parameter and whether it was a named parameter
func (p *ParamSet) FetchMerge(idx int, mergeP Param) (param Param, isNamed bool) {
	name, exists := p.positionToName[idx]
	if !exists || name == "" {
		return mergeP, false
	}

	param, ok := p.namedParams[name]
	if !ok {
		return mergeP, false
	}

	return mergeParam(param, mergeP), true
}

// NewParamSet creates a set of parameters with the given list of already used positions
func NewParamSet(positionsUsed map[int]bool, hasNamedSupport bool) *ParamSet {
	positionToName := make(map[int]string, len(positionsUsed))
	for index, used := range positionsUsed {
		if !used {
			continue
		}

		// assume the previously used params have no name
		positionToName[index] = ""
	}

	return &ParamSet{
		argn:            1,
		namedParams:     make(map[string]Param),
		namedLocs:       make(map[string][]int),
		hasNamedSupport: hasNamedSupport,
		positionToName:  positionToName,
	}
}
