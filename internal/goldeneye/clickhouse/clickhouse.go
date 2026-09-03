// Package clickhouse generates the ClickHouse dialect seed under
// internal/engine/clickhouse/dialect from a clickhouse binary, run as an
// ephemeral `clickhouse local` process that needs no server.
//
// types.jsonl comes from system.data_type_families: every family that is not
// an alias becomes a type, the families that alias it become its aliases,
// and its category is decided by its name, since ClickHouse records no such
// thing. ClickHouse describes its functions no further than their names —
// system.functions carries no signatures — so functions.jsonl is written by
// hand and is not this package's business.
//
// The binary is downloaded once per pinned version by Install, or supplied
// through the CLICKHOUSE environment variable.
package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name of the engine directory the dialect lives under.
const Engine = "clickhouse"

// typeFamilies lists every data type ClickHouse knows and, for a spelling
// such as INT or VARCHAR, the type it stands for.
const typeFamilies = `
SELECT name, alias_to
FROM system.data_type_families
ORDER BY name`

// Version reports the release a binary is.
func Version(ctx context.Context, binary string) (string, error) {
	out, err := exec.CommandContext(ctx, binary, "local", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("clickhouse local --version: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Generate reads the dialect from the binary.
func Generate(ctx context.Context, binary string) (dialect.Files, error) {
	results, err := local{binary: binary}.run(ctx, typeFamilies)
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("clickhouse: expected one result set, got %d", len(results))
	}
	var families []family
	for _, row := range results[0].Data {
		var f family
		if err := json.Unmarshal(row["name"], &f.Name); err != nil {
			return nil, fmt.Errorf("clickhouse: system.data_type_families: %w", err)
		}
		if err := json.Unmarshal(row["alias_to"], &f.AliasTo); err != nil {
			return nil, fmt.Errorf("clickhouse: system.data_type_families: %w", err)
		}
		families = append(families, f)
	}
	types, err := dialect.JSONL(readTypes(families))
	if err != nil {
		return nil, err
	}
	return dialect.Files{dialect.TypesFile: types}, nil
}

type family struct {
	Name    string
	AliasTo string
}

// readTypes turns the families into types: one per family that is not an
// alias, carrying the spellings that alias it. The seed reads names without
// regard to case, so an alias that only differs in case from its type —
// bool for Bool, ENUM for Enum — is dropped.
func readTypes(families []family) []dialect.Type {
	byName := map[string]*dialect.Type{}
	var names []string
	for _, f := range families {
		if f.AliasTo != "" {
			continue
		}
		byName[f.Name] = &dialect.Type{Name: f.Name, Category: category(f.Name)}
		names = append(names, f.Name)
	}
	for _, f := range families {
		if f.AliasTo == "" || strings.EqualFold(f.Name, f.AliasTo) {
			continue
		}
		if t, ok := byName[f.AliasTo]; ok {
			t.Aliases = append(t.Aliases, f.Name)
		}
	}
	sort.Strings(names)
	types := make([]dialect.Type, 0, len(names))
	for _, name := range names {
		types = append(types, *byName[name])
	}
	return types
}

// category assigns a type the PostgreSQL category letter the seed package
// uses, by the family's name: N for numbers, B for booleans, S for strings
// and the identifiers stored as them, D for dates and times, T for time
// spans, A for arrays and U for everything else.
func category(name string) string {
	switch {
	case name == "Bool":
		return "B"
	case name == "Array":
		return "A"
	case strings.HasPrefix(name, "Interval"):
		return "T"
	case hasPrefix(name, "UInt", "Int", "Float", "BFloat", "Decimal"):
		return "N"
	case name == "String", name == "FixedString", name == "UUID", name == "IPv4", name == "IPv6":
		return "S"
	case hasPrefix(name, "Date", "Time"):
		return "D"
	}
	return "U"
}

func hasPrefix(name string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
