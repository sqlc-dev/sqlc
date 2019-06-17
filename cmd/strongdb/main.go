package main

import (
	"flag"
	"log"

	"github.com/kyleconroy/strongdb"
)

func main() {
	pkg := flag.String("package", "db", "package name for Go code")
	sch := flag.String("schema", "", "input directory of SQL migrations")
	out := flag.String("out", "db.go", "output file")
	flag.Parse()

	if err := strongdb.Exec(*sch, flag.Arg(0), *pkg, *out); err != nil {
		log.Fatal(err)
	}
}
