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

// Visitor is the Starlark AST visitor.
type Visitor struct {
	FunctionCalls []*starlark.FunctionCall
}

func (v *Visitor) visit(e build.Expr, stk []build.Expr) {
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
	v.FunctionCalls = append(v.FunctionCalls, fc)
}

func createVisitor() (*Visitor, func(build.Expr, []build.Expr)) {
	v := new(Visitor)
	return v, func(e build.Expr, stk []build.Expr) {
		v.visit(e, stk)
	}
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
	v, visit := createVisitor()
	for _, s := range f.Stmt {
		build.Walk(s, visit)
	}
	log.Printf("Found %d function calls.", len(v.FunctionCalls))
}
