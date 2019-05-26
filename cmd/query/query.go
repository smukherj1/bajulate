package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/bazelbuild/buildtools/build"
	"github.com/smukherj1/bajulate/pkg/starlark"
)

var (
	filePath = flag.String("f", "", "Path to Starlark file to query for bajulation.")
)

func visit(e build.Expr, stk []build.Expr) {
	ce, ok := e.(*build.CallExpr)
	if !ok {
		return
	}

	fc, err := starlark.NewFunctionCall(ce)
	if err != nil {
		log.Printf("Skipping function call on line %d: %v", ce.Pos.Line, err)
		return
	}
	log.Printf("Line %d: %v", fc.Expr.Pos.Line, fc)
}

func main() {
	flag.Parse()
	blob, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read Stalark file %s: %v", *filePath, err)
	}
	f, err := build.Parse(*filePath, blob)
	if err != nil {
		log.Fatalf("Failed to parse Starlark file %s: %v", *filePath, err)
	}
	for _, s := range f.Stmt {
		build.Walk(s, visit)
	}
}
