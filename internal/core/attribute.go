package core

import (
	"database/sql"
	"fmt"
)

// AttributeSpec carries the data needed to register a column on a
// relation. It is the full-fidelity counterpart to CreateAttribute,
// which is kept as a thin wrapper for callers that don't need to set
// any of the "extra" metadata.
type AttributeSpec struct {
	ClassOID       int64
	Name           string
	TypeOID        int64
	Num            int // 1-based ordinal position
	NotNull        bool
	HasDefault     bool
	DeclType       string // original declared type, before normalization
	TypeLength     int    // varchar(N), numeric(p,s).p, etc.
	TypeScale      int    // numeric(p,s).s
	AutoIncrement  bool
	IsPrimaryKey   bool
	IsUnique       bool
}

// CreateAttributeSpec inserts a column with the full set of attribute
// metadata.
func (c *Catalog) CreateAttributeSpec(s AttributeSpec) error {
	_, err := c.db.Exec(`
		INSERT INTO sql_attribute (
			class_oid, name, type_oid, not_null, has_default, num,
			decl_type, type_length, type_scale,
			auto_increment, is_primary_key, is_unique
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ClassOID, s.Name, s.TypeOID,
		boolToInt(s.NotNull), boolToInt(s.HasDefault), s.Num,
		s.DeclType, s.TypeLength, s.TypeScale,
		boolToInt(s.AutoIncrement),
		boolToInt(s.IsPrimaryKey),
		boolToInt(s.IsUnique),
	)
	if err != nil {
		return fmt.Errorf("create attribute %q on class %d: %w", s.Name, s.ClassOID, err)
	}
	return nil
}

// CreateAttribute inserts a column for the given relation. Wraps
// CreateAttributeSpec for callers that only set the basics.
func (c *Catalog) CreateAttribute(classOID int64, name string, typeOID int64, notNull bool, hasDefault bool, num int) error {
	return c.CreateAttributeSpec(AttributeSpec{
		ClassOID:   classOID,
		Name:       name,
		TypeOID:    typeOID,
		Num:        num,
		NotNull:    notNull,
		HasDefault: hasDefault,
	})
}

// SetAttributePrimaryKey flips the is_primary_key (and not_null) flag
// on every named column of a relation. Used when a table-level PRIMARY
// KEY constraint is processed after the columns have been registered.
func (c *Catalog) SetAttributePrimaryKey(classOID int64, columns []string) error {
	for _, name := range columns {
		_, err := c.db.Exec(
			`UPDATE sql_attribute SET is_primary_key = 1, not_null = 1
			 WHERE class_oid = ? AND name = ?`,
			classOID, name,
		)
		if err != nil {
			return fmt.Errorf("mark pk %s on class %d: %w", name, classOID, err)
		}
	}
	return nil
}

// SetAttributeUnique flips is_unique on every named column. Multi-column
// UNIQUE constraints don't make individual columns unique, so callers
// should pass length-1 column lists only.
func (c *Catalog) SetAttributeUnique(classOID int64, columns []string) error {
	for _, name := range columns {
		_, err := c.db.Exec(
			`UPDATE sql_attribute SET is_unique = 1 WHERE class_oid = ? AND name = ?`,
			classOID, name,
		)
		if err != nil {
			return fmt.Errorf("mark unique %s on class %d: %w", name, classOID, err)
		}
	}
	return nil
}

// ColumnInfo describes a resolved column.
type ColumnInfo struct {
	Name            string
	TypeName        string
	TypeOID         int64
	NotNull         bool
	DeclType        string
	TypeLength      int
	TypeScale       int
	AutoIncrement   bool
	IsPrimaryKey    bool
	IsUnique        bool
	AttributeOID    int64
	ClassOID        int64
}

// ResolveColumn looks up a column by table name and column name.
func (c *Catalog) ResolveColumn(table, column string) (*ColumnInfo, error) {
	var info ColumnInfo
	var ai, ipk, iu int
	err := c.db.QueryRow(`
		SELECT a.oid, a.class_oid, a.name, t.name, a.type_oid, a.not_null,
		       a.decl_type, a.type_length, a.type_scale,
		       a.auto_increment, a.is_primary_key, a.is_unique
		FROM sql_attribute a
		JOIN sql_class c ON c.oid = a.class_oid
		JOIN sql_type t ON t.oid = a.type_oid
		WHERE c.name = ? AND a.name = ?
	`, table, column).Scan(
		&info.AttributeOID, &info.ClassOID, &info.Name, &info.TypeName, &info.TypeOID, &info.NotNull,
		&info.DeclType, &info.TypeLength, &info.TypeScale,
		&ai, &ipk, &iu,
	)
	if err != nil {
		return nil, fmt.Errorf("resolve %s.%s: %w", table, column, err)
	}
	info.AutoIncrement = ai != 0
	info.IsPrimaryKey = ipk != 0
	info.IsUnique = iu != 0
	return &info, nil
}

// TableColumns returns all columns for a given table name, ordered by position.
func (c *Catalog) TableColumns(table string) ([]ColumnInfo, error) {
	rows, err := c.db.Query(`
		SELECT a.oid, a.class_oid, a.name, t.name, a.type_oid, a.not_null,
		       a.decl_type, a.type_length, a.type_scale,
		       a.auto_increment, a.is_primary_key, a.is_unique
		FROM sql_attribute a
		JOIN sql_class c ON c.oid = a.class_oid
		JOIN sql_type t ON t.oid = a.type_oid
		WHERE c.name = ?
		ORDER BY a.num
	`, table)
	if err != nil {
		return nil, fmt.Errorf("table columns %q: %w", table, err)
	}
	defer rows.Close()

	var cols []ColumnInfo
	for rows.Next() {
		var ci ColumnInfo
		var ai, ipk, iu int
		if err := rows.Scan(
			&ci.AttributeOID, &ci.ClassOID, &ci.Name, &ci.TypeName, &ci.TypeOID, &ci.NotNull,
			&ci.DeclType, &ci.TypeLength, &ci.TypeScale,
			&ai, &ipk, &iu,
		); err != nil {
			return nil, err
		}
		ci.AutoIncrement = ai != 0
		ci.IsPrimaryKey = ipk != 0
		ci.IsUnique = iu != 0
		cols = append(cols, ci)
	}
	return cols, rows.Err()
}

// AttributeDetails describes a column resolved by its OID. Populated
// from sql_attribute joined back to sql_class / sql_namespace, so the
// caller doesn't have to re-resolve names from raw OIDs.
type AttributeDetails struct {
	Schema          string
	Table           string
	Column          string
	Num             int
	DeclType        string
	TypeLength      int
	TypeScale       int
	AutoIncrement   bool
	IsPrimaryKey    bool
	IsUnique        bool
	NotNull         bool
}

// LookupAttribute returns the full source-side metadata for a column,
// keyed by attribute OID. Used by analyzers to populate Column.Source
// (and the per-column flag fields) once they've resolved an
// expression to a sql_attribute row.
func (c *Catalog) LookupAttribute(attOID int64) (AttributeDetails, error) {
	var ad AttributeDetails
	var schema sql.NullString
	var ai, ipk, iu, nn int
	err := c.db.QueryRow(`
		SELECT ns.name, cls.name, a.name, a.num,
		       a.decl_type, a.type_length, a.type_scale,
		       a.auto_increment, a.is_primary_key, a.is_unique, a.not_null
		FROM sql_attribute a
		JOIN sql_class cls ON cls.oid = a.class_oid
		JOIN sql_namespace ns ON ns.oid = cls.namespace_oid
		WHERE a.oid = ?
	`, attOID).Scan(
		&schema, &ad.Table, &ad.Column, &ad.Num,
		&ad.DeclType, &ad.TypeLength, &ad.TypeScale,
		&ai, &ipk, &iu, &nn,
	)
	if err != nil {
		return ad, fmt.Errorf("lookup attribute %d: %w", attOID, err)
	}
	ad.Schema = schema.String
	ad.AutoIncrement = ai != 0
	ad.IsPrimaryKey = ipk != 0
	ad.IsUnique = iu != 0
	ad.NotNull = nn != 0
	return ad, nil
}
