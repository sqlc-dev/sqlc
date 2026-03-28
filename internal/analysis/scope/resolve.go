package scope

import "fmt"

// ResolutionError describes why name resolution failed with full provenance.
type ResolutionError struct {
	Name       string
	Qualifier  string // Table/alias qualifier, if any
	Kind       ResolutionErrorKind
	Scope      *Scope   // The scope where resolution was attempted
	Candidates []string // For ambiguity errors, the competing names
	Location   int      // Source position of the reference
}

type ResolutionErrorKind int

const (
	ErrNotFound ResolutionErrorKind = iota
	ErrAmbiguous
	ErrQualifierNotFound // e.g., "u.name" but "u" doesn't exist
)

func (e *ResolutionError) Error() string {
	switch e.Kind {
	case ErrNotFound:
		if e.Qualifier != "" {
			return fmt.Sprintf("column %q not found in %q", e.Name, e.Qualifier)
		}
		return fmt.Sprintf("column %q does not exist", e.Name)
	case ErrAmbiguous:
		return fmt.Sprintf("column reference %q is ambiguous", e.Name)
	case ErrQualifierNotFound:
		return fmt.Sprintf("table or alias %q does not exist", e.Qualifier)
	default:
		return fmt.Sprintf("resolution error for %q", e.Name)
	}
}

// ResolutionPath records the edges traversed during successful resolution.
// This is the provenance — it tells you exactly how a name was resolved.
type ResolutionPath struct {
	Steps []ResolutionStep
}

type ResolutionStep struct {
	Edge  *Edge  // nil for the final lookup step
	Scope *Scope // The scope where this step occurred
}

// ResolvedName is the result of successful name resolution.
type ResolvedName struct {
	Declaration *Declaration
	Path        ResolutionPath
}

// Resolve looks up an unqualified column name in this scope.
// It searches local declarations first, then follows parent edges.
// Returns an error if the name is not found or is ambiguous.
func (s *Scope) Resolve(name string) (*ResolvedName, error) {
	return s.resolve(name, nil, 0)
}

// ResolveQualified looks up a qualified name like "u.name".
// First resolves the qualifier (table/alias), then looks up the column
// in that table's scope.
func (s *Scope) ResolveQualified(qualifier, name string) (*ResolvedName, error) {
	// First, find the qualifier (table name or alias)
	qualScope, err := s.resolveQualifier(qualifier, 0)
	if err != nil {
		return nil, &ResolutionError{
			Name:      name,
			Qualifier: qualifier,
			Kind:      ErrQualifierNotFound,
			Scope:     s,
		}
	}

	// Then resolve the column within that scope
	var matches []*Declaration
	for _, d := range qualScope.Declarations {
		if d.Name == name {
			matches = append(matches, d)
		}
	}

	if len(matches) == 0 {
		return nil, &ResolutionError{
			Name:      name,
			Qualifier: qualifier,
			Kind:      ErrNotFound,
			Scope:     qualScope,
		}
	}
	if len(matches) > 1 {
		return nil, &ResolutionError{
			Name:      name,
			Qualifier: qualifier,
			Kind:      ErrAmbiguous,
			Scope:     qualScope,
		}
	}

	return &ResolvedName{
		Declaration: matches[0],
		Path: ResolutionPath{
			Steps: []ResolutionStep{
				{Scope: s},
				{Scope: qualScope},
			},
		},
	}, nil
}

const maxResolutionDepth = 20

// resolve performs recursive name resolution with cycle detection via depth limit.
func (s *Scope) resolve(name string, visited map[*Scope]bool, depth int) (*ResolvedName, error) {
	if depth > maxResolutionDepth {
		return nil, &ResolutionError{Name: name, Kind: ErrNotFound, Scope: s}
	}
	if visited == nil {
		visited = make(map[*Scope]bool)
	}
	if visited[s] {
		return nil, &ResolutionError{Name: name, Kind: ErrNotFound, Scope: s}
	}
	visited[s] = true

	// Search local declarations first
	var matches []*Declaration
	for _, d := range s.Declarations {
		if d.Name == name && d.Kind == DeclColumn {
			matches = append(matches, d)
		}
	}

	// Also search table/alias declarations to find columns inside their scopes
	for _, d := range s.Declarations {
		if (d.Kind == DeclTable || d.Kind == DeclAlias || d.Kind == DeclCTE) && d.Scope != nil {
			for _, cd := range d.Scope.Declarations {
				if cd.Name == name && cd.Kind == DeclColumn {
					matches = append(matches, cd)
				}
			}
		}
	}

	if len(matches) == 1 {
		return &ResolvedName{
			Declaration: matches[0],
			Path: ResolutionPath{
				Steps: []ResolutionStep{{Scope: s}},
			},
		}, nil
	}
	if len(matches) > 1 {
		return nil, &ResolutionError{Name: name, Kind: ErrAmbiguous, Scope: s}
	}

	// Follow parent, lateral, and outer edges
	for _, edge := range s.Edges {
		switch edge.Kind {
		case EdgeParent, EdgeLateral, EdgeOuter:
			result, err := edge.Target.resolve(name, visited, depth+1)
			if err == nil {
				result.Path.Steps = append([]ResolutionStep{{Edge: edge, Scope: s}}, result.Path.Steps...)
				return result, nil
			}
			// Propagate ambiguity errors — don't swallow them
			if resErr, ok := err.(*ResolutionError); ok && resErr.Kind == ErrAmbiguous {
				return nil, resErr
			}
		}
	}

	return nil, &ResolutionError{Name: name, Kind: ErrNotFound, Scope: s}
}

// resolveQualifier finds the scope associated with a table name or alias.
func (s *Scope) resolveQualifier(qualifier string, depth int) (*Scope, error) {
	if depth > maxResolutionDepth {
		return nil, fmt.Errorf("qualifier %q not found", qualifier)
	}

	// Check alias edges first (higher priority)
	for _, edge := range s.Edges {
		if edge.Kind == EdgeAlias && edge.Label == qualifier {
			return edge.Target, nil
		}
	}

	// Check local table/alias declarations
	for _, d := range s.Declarations {
		if d.Name == qualifier && (d.Kind == DeclTable || d.Kind == DeclAlias || d.Kind == DeclCTE) && d.Scope != nil {
			return d.Scope, nil
		}
	}

	// Follow parent edges
	for _, edge := range s.Edges {
		if edge.Kind == EdgeParent || edge.Kind == EdgeLateral || edge.Kind == EdgeOuter {
			result, err := edge.Target.resolveQualifier(qualifier, depth+1)
			if err == nil {
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("qualifier %q not found", qualifier)
}

// ResolveColumnRef resolves a column reference that may have 1, 2, or 3 parts:
//   - ["name"]           -> unqualified column
//   - ["alias", "name"]  -> table-qualified column
//   - ["schema", "table", "name"] -> schema-qualified column (treated as qualifier="table")
func (s *Scope) ResolveColumnRef(parts []string) (*ResolvedName, error) {
	switch len(parts) {
	case 1:
		return s.Resolve(parts[0])
	case 2:
		return s.ResolveQualified(parts[0], parts[1])
	case 3:
		// For now, ignore schema and use table.column
		return s.ResolveQualified(parts[1], parts[2])
	default:
		return nil, fmt.Errorf("invalid column reference with %d parts", len(parts))
	}
}

// AllColumns returns all column declarations visible from this scope,
// optionally filtered by a qualifier. This is used for SELECT * expansion.
func (s *Scope) AllColumns(qualifier string) []*Declaration {
	if qualifier != "" {
		qualScope, err := s.resolveQualifier(qualifier, 0)
		if err != nil {
			return nil
		}
		var cols []*Declaration
		for _, d := range qualScope.Declarations {
			if d.Kind == DeclColumn {
				cols = append(cols, d)
			}
		}
		return cols
	}

	// Collect from all table/alias declarations in this scope
	var cols []*Declaration
	seen := make(map[string]bool)

	var collect func(sc *Scope, depth int)
	collect = func(sc *Scope, depth int) {
		if depth > maxResolutionDepth {
			return
		}
		for _, d := range sc.Declarations {
			if (d.Kind == DeclTable || d.Kind == DeclAlias || d.Kind == DeclCTE) && d.Scope != nil {
				for _, cd := range d.Scope.Declarations {
					if cd.Kind == DeclColumn && !seen[d.Name+"."+cd.Name] {
						seen[d.Name+"."+cd.Name] = true
						cols = append(cols, cd)
					}
				}
			}
		}
		for _, edge := range sc.Edges {
			if edge.Kind == EdgeParent {
				collect(edge.Target, depth+1)
			}
		}
	}
	collect(s, 0)
	return cols
}
