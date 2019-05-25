package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/bazelbuild/buildtools/build"
)

var (
	filePath = flag.String("f", "", "Path to Starlark file to query for bajulation.")
)

func visit(e build.Expr, stk []build.Expr) {
	//log.Printf("Visit %T.\n", e)
	ce, ok := e.(*build.CallExpr)
	if !ok {
		return
	}
	id, ok := ce.X.(*build.Ident)
	if !ok {
		return
	}
	s, _ := ce.Span()
	log.Printf("%s call on line %d.\n", id.Name, s.Line)
}

func main() {
	flag.Parse()
	log.Println("Running the Bajulate Query Utility.")
	blob, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Unable to read Stalark file %s: %v", *filePath, err)
	}
	f, err := build.Parse(*filePath, blob)
	if err != nil {
		log.Fatalf("Unable to parse Starlark file %s: %v", *filePath, err)
	}
	for _, s := range f.Stmt {
		build.Walk(s, visit)
	}
	log.Println("Bajulate Query Utility is done.")
}
