package starlark

import (
	"fmt"
	"strings"

	"github.com/bazelbuild/buildtools/build"
	"github.com/pkg/errors"
)

// FunctionArg is a wrapper around a Starlark Function argument.
type FunctionArg struct {
	ParamExpr         *build.Ident
	SingleValueExpr   *build.StringExpr
	MultipleValueExpr *build.ListExpr
}

// Param returns the name of the parameter in this function argument refers to.
func (a *FunctionArg) Param() string {
	return a.ParamExpr.Name
}

// HasSingleValue eturns whether this function argument was assigned a single
// string value in the call expression.
func (a *FunctionArg) HasSingleValue() bool {
	return a.SingleValueExpr != nil
}

// SingleValue returns the string value assigned to this function argument.
func (a *FunctionArg) SingleValue() string {
	return a.SingleValueExpr.Value
}

// HasMultipleValue returns whether this function argument was assigned a list
// of strings as value.
func (a *FunctionArg) HasMultipleValue() bool {
	return a.MultipleValueExpr != nil
}

// MultipleValue returns the list of strings that was assigned to this function
// argument.
func (a *FunctionArg) MultipleValue() []string {
	result := []string{}
	for _, e := range a.MultipleValueExpr.List {
		s := e.(*build.StringExpr)
		result = append(result, s.Value)
	}
	return result
}

func (a *FunctionArg) String() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "%s = ", a.Param())
	if a.HasSingleValue() {
		fmt.Fprintf(b, "%s", a.SingleValue())
	}
	if a.HasMultipleValue() {
		fmt.Fprintf(b, "%v", a.MultipleValue())
	}
	return b.String()
}

// NewFunctionArg creates a function argument wrapper object for the given
// Starlark call expression.
func NewFunctionArg(a *build.AssignExpr) (*FunctionArg, error) {
	result := new(FunctionArg)
	id, ok := a.LHS.(*build.Ident)
	if !ok {
		return nil, fmt.Errorf("LHS of expression was not an identifier")
	}
	result.ParamExpr = id
	s, ok := a.RHS.(*build.StringExpr)
	if ok {
		result.SingleValueExpr = s
		return result, nil
	}
	l, ok := a.RHS.(*build.ListExpr)
	if !ok {
		return nil, fmt.Errorf("unsupported type on RHS of expression %T", a.RHS)
	}
	for _, s := range l.List {
		_, ok := s.(*build.StringExpr)
		if !ok {
			return nil, fmt.Errorf("unsupported element of type %T in list expression in function argument", s)
		}
	}
	result.MultipleValueExpr = l
	return result, nil
}

// FunctionCall is a wrapper around a Starlark Function call.
type FunctionCall struct {
	Expr     *build.CallExpr
	NameExpr *build.Ident

	Args []*FunctionArg
}

// Name is the name of the Starlark function.
func (c *FunctionCall) Name() string {
	return c.NameExpr.Name
}

func (c *FunctionCall) String() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "Function %s, line %d, %d arguments: %v\n", c.Name(), c.Expr.Pos.Line, len(c.Args), c.Args)
	return b.String()
}

// NewFunctionCall creates a function call wrapper object for the given
// Starlark call expression.
func NewFunctionCall(c *build.CallExpr) (*FunctionCall, error) {
	result := new(FunctionCall)
	result.Expr = c
	id, ok := c.X.(*build.Ident)
	if !ok {
		return nil, fmt.Errorf("function name is not an identifier")
	}
	result.NameExpr = id
	for _, arg := range c.List {
		asgn, ok := arg.(*build.AssignExpr)
		if !ok {
			return nil, fmt.Errorf("one of the argument expressions on %s function call was not an assignment", id.Name)
		}
		a, err := NewFunctionArg(asgn)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to process argument for %s function call", id.Name)
		}
		result.Args = append(result.Args, a)
	}
	return result, nil
}
