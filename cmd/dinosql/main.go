package main

import (
	"flag"
	"log"

	"github.com/kyleconroy/dinosql"
)

func main() {
	flag.Parse()
	if err := dinosql.Exec(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
}
