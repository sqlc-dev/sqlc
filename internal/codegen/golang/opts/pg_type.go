package opts

var pgTypeCanonicalNames map[string]string

func init() {
	groups := []struct {
		canonical string
		aliases   []string
	}{
		{"pg_catalog.timestamptz", []string{"timestamptz", "timestamp with time zone"}},
		{"pg_catalog.timestamp", []string{"timestamp", "timestamp without time zone"}},
		{"pg_catalog.time", []string{"time", "time without time zone"}},
		{"pg_catalog.timetz", []string{"timetz", "time with time zone"}},
	}

	pgTypeCanonicalNames = make(map[string]string, len(groups)*3)
	for _, g := range groups {
		pgTypeCanonicalNames[g.canonical] = g.canonical
		for _, alias := range g.aliases {
			pgTypeCanonicalNames[alias] = g.canonical
		}
	}
}

// canonicalPostgreSQLType maps PostgreSQL type aliases to a single canonical name
// so db_type overrides match regardless of spelling in schema SQL or config.
func canonicalPostgreSQLType(t string) string {
	if canonical, ok := pgTypeCanonicalNames[t]; ok {
		return canonical
	}
	return t
}
