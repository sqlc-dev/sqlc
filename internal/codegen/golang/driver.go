package golang

type Driver string

const (
	PgxDriver    Driver = "pgx/v4"
	StdLibDriver Driver = "database/sql"
)

func DriverFromString(s string) Driver {
	switch s {
	case string(PgxDriver):
		return PgxDriver
	default:
		return StdLibDriver
	}
}
