package evaluation

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
)

/* Structs */
type (
	/* Collections and collection-likes */

	Tuple struct {
		Value []Object
	}

	Array struct {
		Value []Object
	}

	Map struct {
		Values map[string]Object
		Keys   map[string]Object
	}

	/* Others */

	Block struct {
		Params []ast.Expression
		Body   ast.Statement
	}

	Class struct {
		Name    string
		Parent  Object
		Methods []Object
	}

	Instance struct {
		Base Object
		Data map[string]Object
	}
)

/* Type() methods */
func (_ *Tuple) Type() Type    { return TUPLE }
func (_ *Array) Type() Type    { return ARRAY }
func (_ *Map) Type() Type      { return MAP }
func (_ *Block) Type() Type    { return BLOCK }
func (_ *Class) Type() Type    { return CLASS }
func (_ *Instance) Type() Type { return INSTANCE }

/* Equals() methods */
func (t *Tuple) Equals(o Object) bool {
	if other, ok := o.(*Tuple); ok {
		if len(other.Value) != len(t.Value) {
			return false
		}

		for i, elem := range t.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (a *Array) Equals(o Object) bool {
	if other, ok := o.(*Array); ok {
		if len(other.Value) != len(a.Value) {
			return false
		}

		for i, elem := range a.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (m *Map) Equals(o Object) bool {
	if other, ok := o.(*Map); ok {
		if len(other.Values) != len(m.Values) {
			return false
		}

		for k, v := range m.Values {
			if _, ok := other.Values[k]; !ok {
				return false
			}

			if !v.Equals(other.Values[k]) {
				return false
			}
		}

		return true
	}

	return false
}

func (_ *Block) Equals(o Object) bool {
	_, ok := o.(*Block)
	return ok
}

func (c *Class) Equals(o Object) bool {
	if other, ok := o.(*Class); ok {
		return other.Name == c.Name
	}

	return false
}

func (i *Instance) Equals(o Object) bool {
	if other, ok := o.(*Instance); ok {
		if !other.Base.Equals(i.Base) {
			return false
		}

		for k, v := range i.Data {
			if !v.Equals(other.Data[k]) {
				return false
			}
		}

		return true
	}

	return false
}

/* Stringer implementations */
func join(arr []Object, sep string) string {
	stringArr := make([]string, len(arr))

	for i, elem := range arr {
		stringArr[i] = elem.String()
	}

	return strings.Join(stringArr, ", ")
}

func (t *Tuple) String() string {
	return fmt.Sprintf("(%s)", join(t.Value, ", "))
}

func (a *Array) String() string {
	return fmt.Sprintf("[%s]", join(a.Value, ", "))
}

func (m *Map) String() string {
	stringArr := make([]string, len(m.Values))
	i := 0

	for k, v := range m.Values {
		stringArr[i] = fmt.Sprintf(
			"%s: %s",
			m.Keys[k].String(),
			v.String(),
		)

		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(stringArr, ", "))
}

func (b *Block) String() string {
	return "<block>"
}

func (c *Class) String() string {
	return c.Name
}

func (i *Instance) String() string {
	if i.Base.(*Class).Name == "Error" {
		return fmt.Sprintf("%s: %s", i.Get(&String{"tag"}), i.Get(&String{"msg"}))
	}

	stringMethod := i.Base.(*Class).GetMethod("string")

	if stringMethod != nil {
		args := map[string]Object{
			"self": i,
		}

		enclosed := stringMethod.Fn.Context.EncloseWith(args)
		result := eval(stringMethod.Fn.Body, enclosed)

		return result.String()
	}

	return fmt.Sprintf("<instance of %s>", i.Base.(*Class).Name)
}

/* Collection implementations */
func (t *Tuple) Elements() []Object {
	return t.Value
}

func (t *Tuple) GetIndex(i int) Object {
	if i >= len(t.Value) || i < 0 {
		return O_NULL
	}

	return t.Value[i]
}

func (t *Tuple) SetIndex(i int, o Object) {
	if i >= len(t.Value) || i < 0 {
		return
	}

	t.Value[i] = o
}

func (a *Array) Elements() []Object {
	return a.Value
}

func (a *Array) GetIndex(i int) Object {
	if i >= len(a.Value) || i < 0 {
		return O_NULL
	}

	return a.Value[i]
}

func (a *Array) SetIndex(i int, o Object) {
	if i >= len(a.Value) || i < 0 {
		return
	}

	a.Value[i] = o
}

/* Container implementations */
func (m *Map) Get(key Object) Object {
	if hasher, ok := key.(Hasher); !ok {
		return O_NULL
	} else {
		if val, ok := m.Values[hasher.Hash()]; ok {
			return val
		}

		return O_NULL
	}
}

func (m *Map) Set(key, value Object) {
	if hasher, ok := key.(Hasher); ok {
		hash := hasher.Hash()
		m.Values[hash] = value
		m.Keys[hash] = key
	}
}

func (i *Instance) Get(key Object) Object {
	if strkey, ok := key.(*String); !ok {
		return O_NULL
	} else {
		if val, ok := i.Data[strkey.Value]; ok {
			return val
		}

		return O_NULL
	}
}

func (i *Instance) Set(key, value Object) {
	if strkey, ok := key.(*String); ok {
		i.Data[strkey.Value] = value
	}
}

/* Other methods */
func (c *Class) GetMethods() []Method {
	var methods []Method

	if c.Parent != nil {
		methods = c.Parent.(*Class).GetMethods()
	}

	for _, m := range c.Methods {
		if method, ok := m.(*Method); ok {
			methods = append(methods, *method)
		}
	}

	return methods
}

func (c *Class) GetMethod(pattern string) *Method {
	fnPattern := strings.Split(pattern, " ")

	for _, method := range c.GetMethods() {
		methodPattern := method.Fn.Pattern

		if len(fnPattern) != len(methodPattern) {
			continue
		}

		isMatch := true
		for i, mPatItem := range methodPattern {
			_, isParam := mPatItem.(*ast.Parameter)
			_, isIdent := mPatItem.(*ast.Identifier)

			if !(fnPattern[i] == "$" && isParam || isIdent && fnPattern[i] == methodPattern[i].Token().Literal) {
				isMatch = false
			}
		}

		if isMatch {
			return &method
		}
	}

	return nil
}
