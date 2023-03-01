package golang

type Field struct {
	Name    string // CamelCased name for Go
	DBName  string // Name as used in the DB
	Type    string
	Tags    map[string]string
	Comment string
}

func (gf Field) Tag() string {
	return TagsToString(gf.Tags)
}
