package evaluation

import "fmt"

// Object is the base object interface, which
// every object implements
type Object interface {
	fmt.Stringer
	Equals(Object) bool
	Type() Type
}

// Collection is a child interface of Object,
// which represents an object which can be
// thought of as a list of items
type Collection interface {
	Object
	Elements() []Object
	GetIndex(int) Object
	SetIndex(int, Object)
}

// Container is a child interface of Object,
// which can be accessed by keys - like a map
type Container interface {
	Object
	Get(Object) Object
	Set(Object, Object)
}

// Hasher is any object which can be a key
// in a map
type Hasher interface {
	Object
	Hash() string
}
