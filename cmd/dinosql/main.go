package main

import (
	"flag"
	"log"

	"github.com/kyleconroy/dinosql"
)

func main() {
	pkg := flag.String("package", "db", "package name for Go code")
	sch := flag.String("schema", "", "input directory of SQL migrations")
	prepare := flag.Bool("prepare", false, "include prepared query support")
	tags := flag.Bool("tags", false, "add tags to database records")
	out := flag.String("out", "db.go", "output file")
	flag.Parse()

	settings := dinosql.GenerateSettings{
		Package:             *pkg,
		EmitPreparedQueries: *prepare,
		EmitTags:            *tags,
	}

	if err := dinosql.Exec(*sch, flag.Arg(0), *out, settings); err != nil {
		log.Fatal(err)
	}
}
