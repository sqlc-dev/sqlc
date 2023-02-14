package golang

type SQLDriver int

const (
	SQLPackagePGXV4    string = "pgx/v4"
	SQLPackagePGXV5    string = "pgx/v5"
	SQLPackageStandard string = "database/sql"
)

const (
	SQLDriverPGXV4 SQLDriver = iota
	SQLDriverPGXV5
	SQLDriverLibPQ
)

func parseDriver(sqlPackage string) SQLDriver {
	switch sqlPackage {
	case SQLPackagePGXV4:
		return SQLDriverPGXV4
	case SQLPackagePGXV5:
		return SQLDriverPGXV5
	default:
		return SQLDriverLibPQ
	}
}

func (d SQLDriver) IsPGX() bool {
	return d == SQLDriverPGXV4 || d == SQLDriverPGXV5
}

func (d SQLDriver) Package() string {
	switch d {
	case SQLDriverPGXV4:
		return SQLPackagePGXV4
	case SQLDriverPGXV5:
		return SQLPackagePGXV5
	default:
		return SQLPackageStandard
	}
}
