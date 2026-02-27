// Package testcases defines the standard set of end-to-end test cases
// that each SQL engine must implement.
//
// Each engine (PostgreSQL, MySQL, SQLite) provides its own implementation
// of these test cases with engine-specific SQL syntax. The registry only
// defines WHAT to test, not HOW - each engine's testdata directory contains
// the actual SQL files.
//
// Directory structure for each engine:
//
//	internal/endtoend/{engine}/testdata/
//	├── core/                    # Core tests (required)
//	│   ├── select_star/
//	│   │   ├── sqlc.yaml
//	│   │   ├── query.sql
//	│   │   └── go/              # Expected output
//	│   └── ...
//	├── enum/                    # Enum extension (optional)
//	│   └── ...
//	├── schema/                  # Schema namespace extension (optional)
//	│   └── ...
//	└── schema.sql               # Shared schema for this engine
package testcases

// Category represents a category of test cases
type Category string

const (
	// Core categories - all engines must implement these
	CategorySelect    Category = "select"
	CategoryInsert    Category = "insert"
	CategoryUpdate    Category = "update"
	CategoryDelete    Category = "delete"
	CategoryJoin      Category = "join"
	CategoryCTE       Category = "cte"
	CategorySubquery  Category = "subquery"
	CategoryUnion     Category = "union"
	CategoryAggregate Category = "aggregate"
	CategoryOperator  Category = "operator"
	CategoryCase      Category = "case"
	CategoryNull      Category = "null"
	CategoryCast      Category = "cast"
	CategoryFunction  Category = "function"
	CategoryDataType  Category = "datatype"
	CategoryDDL       Category = "ddl"
	CategoryView      Category = "view"
	CategoryUpsert    Category = "upsert"
	CategoryParam     Category = "param"
	CategoryResult    Category = "result"
	CategoryError     Category = "error"

	// Extension categories - optional based on engine capabilities
	CategoryEnum   Category = "enum"
	CategorySchema Category = "schema"
	CategoryArray  Category = "array"
	CategoryJSON   Category = "json"
)

// TestCase defines a single end-to-end test case
type TestCase struct {
	// ID is the unique identifier for this test case (e.g., "S01")
	ID string

	// Name is the test case name used in the filesystem (e.g., "select_star")
	Name string

	// Category is the category this test belongs to
	Category Category

	// Description explains what this test validates
	Description string

	// Required indicates if this test is mandatory for all engines
	Required bool
}

// Registry holds all registered test cases
type Registry struct {
	cases    map[string]*TestCase
	byCategory map[Category][]*TestCase
}

// NewRegistry creates a new test case registry
func NewRegistry() *Registry {
	return &Registry{
		cases:      make(map[string]*TestCase),
		byCategory: make(map[Category][]*TestCase),
	}
}

// Register adds a test case to the registry
func (r *Registry) Register(tc *TestCase) {
	r.cases[tc.ID] = tc
	r.byCategory[tc.Category] = append(r.byCategory[tc.Category], tc)
}

// Get returns a test case by ID
func (r *Registry) Get(id string) *TestCase {
	return r.cases[id]
}

// GetByCategory returns all test cases in a category
func (r *Registry) GetByCategory(cat Category) []*TestCase {
	return r.byCategory[cat]
}

// All returns all registered test cases
func (r *Registry) All() []*TestCase {
	result := make([]*TestCase, 0, len(r.cases))
	for _, tc := range r.cases {
		result = append(result, tc)
	}
	return result
}

// Required returns all required test cases
func (r *Registry) Required() []*TestCase {
	var result []*TestCase
	for _, tc := range r.cases {
		if tc.Required {
			result = append(result, tc)
		}
	}
	return result
}

// RequiredCategories returns the categories that all engines must implement
func RequiredCategories() []Category {
	return []Category{
		CategorySelect,
		CategoryInsert,
		CategoryUpdate,
		CategoryDelete,
		CategoryJoin,
		CategoryCTE,
		CategorySubquery,
		CategoryUnion,
		CategoryAggregate,
		CategoryOperator,
		CategoryCase,
		CategoryNull,
		CategoryCast,
		CategoryFunction,
		CategoryDataType,
		CategoryDDL,
		CategoryView,
		CategoryUpsert,
		CategoryParam,
		CategoryResult,
		CategoryError,
	}
}

// ExtensionCategories returns optional extension categories
func ExtensionCategories() []Category {
	return []Category{
		CategoryEnum,
		CategorySchema,
		CategoryArray,
		CategoryJSON,
	}
}

// Engine represents a SQL database engine
type Engine string

const (
	EnginePostgreSQL Engine = "postgresql"
	EngineMySQL      Engine = "mysql"
	EngineSQLite     Engine = "sqlite"
)

// EngineCapabilities defines what features an engine supports
type EngineCapabilities struct {
	// SupportsReturning indicates if the engine supports RETURNING clause
	SupportsReturning bool

	// SupportsFullOuterJoin indicates if the engine supports FULL OUTER JOIN
	SupportsFullOuterJoin bool

	// SupportsRightJoin indicates if the engine supports RIGHT JOIN
	SupportsRightJoin bool

	// SupportsCTE indicates if the engine supports Common Table Expressions
	SupportsCTE bool

	// SupportsRecursiveCTE indicates if the engine supports recursive CTEs
	SupportsRecursiveCTE bool

	// SupportsUpsert indicates if the engine supports upsert operations
	SupportsUpsert bool

	// SupportsEnum indicates if the engine supports ENUM types
	SupportsEnum bool

	// SupportsSchema indicates if the engine supports schema namespaces
	SupportsSchema bool

	// SupportsArray indicates if the engine supports array types
	SupportsArray bool

	// SupportsJSON indicates if the engine supports native JSON types
	SupportsJSON bool

	// SupportsIntersect indicates if the engine supports INTERSECT
	SupportsIntersect bool

	// SupportsExcept indicates if the engine supports EXCEPT
	SupportsExcept bool
}

// DefaultCapabilities returns the default capabilities for each engine
func DefaultCapabilities(engine Engine) EngineCapabilities {
	switch engine {
	case EnginePostgreSQL:
		return EngineCapabilities{
			SupportsReturning:     true,
			SupportsFullOuterJoin: true,
			SupportsRightJoin:     true,
			SupportsCTE:           true,
			SupportsRecursiveCTE:  true,
			SupportsUpsert:        true,
			SupportsEnum:          true,
			SupportsSchema:        true,
			SupportsArray:         true,
			SupportsJSON:          true,
			SupportsIntersect:     true,
			SupportsExcept:        true,
		}
	case EngineMySQL:
		return EngineCapabilities{
			SupportsReturning:     false, // MySQL 8.0.21+ has limited support
			SupportsFullOuterJoin: false,
			SupportsRightJoin:     true,
			SupportsCTE:           true, // MySQL 8.0+
			SupportsRecursiveCTE:  true, // MySQL 8.0+
			SupportsUpsert:        true, // ON DUPLICATE KEY UPDATE
			SupportsEnum:          true,
			SupportsSchema:        true, // databases act as schemas
			SupportsArray:         false,
			SupportsJSON:          true,
			SupportsIntersect:     true, // MySQL 8.0.31+
			SupportsExcept:        true, // MySQL 8.0.31+
		}
	case EngineSQLite:
		return EngineCapabilities{
			SupportsReturning:     true, // SQLite 3.35+
			SupportsFullOuterJoin: false,
			SupportsRightJoin:     true, // SQLite 3.39+
			SupportsCTE:           true,
			SupportsRecursiveCTE:  true,
			SupportsUpsert:        true, // ON CONFLICT
			SupportsEnum:          false,
			SupportsSchema:        false, // attached databases only
			SupportsArray:         false,
			SupportsJSON:          true, // json1 extension
			SupportsIntersect:     true,
			SupportsExcept:        true,
		}
	default:
		return EngineCapabilities{}
	}
}

// TestsForEngine returns all test cases that an engine should implement
// based on its capabilities
func (r *Registry) TestsForEngine(engine Engine) []*TestCase {
	caps := DefaultCapabilities(engine)
	var result []*TestCase

	for _, tc := range r.cases {
		if shouldIncludeTest(tc, caps) {
			result = append(result, tc)
		}
	}
	return result
}

// RequiredTestsForEngine returns required test cases for an engine
func (r *Registry) RequiredTestsForEngine(engine Engine) []*TestCase {
	caps := DefaultCapabilities(engine)
	var result []*TestCase

	for _, tc := range r.cases {
		if tc.Required && shouldIncludeTest(tc, caps) {
			result = append(result, tc)
		}
	}
	return result
}

func shouldIncludeTest(tc *TestCase, caps EngineCapabilities) bool {
	// Extension categories depend on capabilities
	switch tc.Category {
	case CategoryEnum:
		return caps.SupportsEnum
	case CategorySchema:
		return caps.SupportsSchema
	case CategoryArray:
		return caps.SupportsArray
	case CategoryJSON:
		return caps.SupportsJSON
	}

	// Some specific tests depend on capabilities
	switch tc.ID {
	case "J03": // join_right
		return caps.SupportsRightJoin
	case "J04": // join_full
		return caps.SupportsFullOuterJoin
	case "I04", "I05", "U05", "U06", "D03", "D04": // RETURNING tests
		return caps.SupportsReturning
	case "C04": // cte_recursive
		return caps.SupportsRecursiveCTE
	case "C01", "C02", "C03", "C05", "C06", "C07": // CTE tests
		return caps.SupportsCTE
	case "UP01", "UP02", "UP03": // upsert tests
		return caps.SupportsUpsert
	case "N04": // intersect
		return caps.SupportsIntersect
	case "N05": // except
		return caps.SupportsExcept
	}

	return true
}
