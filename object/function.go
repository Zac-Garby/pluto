package object

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
)

/* Structs */
type (
	// Function is a normal Pluto function, referenced by its pattern
	Function struct {
		Pattern   []ast.Expression
		Body      bytecode.Code
		Constants []Object
		Names     []string
		Patterns  []string
		OnCall    func(self *Function) Object
	}
)

/* Type() methods */

// Type returns the type of the object
func (f *Function) Type() Type { return FunctionType }

/* Equals() methods */

// Equals checks if two objects are equal to each other
func (f *Function) Equals(o Object) bool {
	_, ok := o.(*Function)
	return ok
}

/* Stringer implementations */
func (f *Function) String() string {
	var pstring []string

	for _, item := range f.Pattern {
		if _, ok := item.(*ast.Parameter); ok {
			pstring = append(pstring, "$")
		} else {
			pstring = append(pstring, item.Token().Literal)
		}
	}

	return fmt.Sprintf("<function: %s>", strings.Join(pstring, " "))
}