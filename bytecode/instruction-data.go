package bytecode

type data struct {
	Name   string
	HasArg bool
}

// Instructions stores data about the different instruction types
var Instructions = map[byte]data{
	Pop: {Name: "POP"},
	Dup: {Name: "DUP"},
	Rot: {Name: "ROT"},

	LoadConst:  {Name: "LOAD_CONST", HasArg: true},
	LoadName:   {Name: "LOAD_NAME", HasArg: true},
	StoreName:  {Name: "STORE_NAME", HasArg: true},
	LoadField:  {Name: "LOAD_FIELD"},
	StoreField: {Name: "STORE_FIELD"},

	UnaryInvert: {Name: "UNARY_INVERT"},
	UnaryNegate: {Name: "UNARY_NEGATE"},
	UnaryNoOp:   {Name: "UNARY_NO_OP"},

	BinaryAdd:      {Name: "BINARY_ADD"},
	BinarySubtract: {Name: "BINARY_SUBTRACT"},
	BinaryMultiply: {Name: "BINARY_MULTIPLY"},
	BinaryDivide:   {Name: "BINARY_DIVIDE"},
	BinaryExponent: {Name: "BINARY_EXPONENT"},
	BinaryFloorDiv: {Name: "BINARY_FLOOR_DIV"},
	BinaryMod:      {Name: "BINARY_MOD"},
	BinaryOr:       {Name: "BINARY_OR"},
	BinaryAnd:      {Name: "BINARY_AND"},
	BinaryBitOr:    {Name: "BINARY_BIT_OR"},
	BinaryBitAnd:   {Name: "BINARY_BIT_AND"},
	BinaryEquals:   {Name: "BINARY_EQUALS"},
	BinaryNotEqual: {Name: "BINARY_NOT_EQUAL"},
	BinaryLessThan: {Name: "BINARY_LESS_THAN"},
	BinaryMoreThan: {Name: "BINARY_MORE_THAN"},
	BinaryLessEq:   {Name: "BINARY_LESS_EQ"},
	BinaryMoreEq:   {Name: "BINARY_MORE_EQ"},

	PushFn:  {Name: "PUSH_FN", HasArg: true},
	CallFn:  {Name: "CALL_FN"},
	Return:  {Name: "RETURN_FN"},
	DoBlock: {Name: "DO_BLOCK"},

	Print:   {Name: "PRINT"},
	Println: {Name: "PRINT_LINE"},

	Jump:        {Name: "JUMP", HasArg: true},
	JumpIfTrue:  {Name: "JUMP_IF_TRUE", HasArg: true},
	JumpIfFalse: {Name: "JUMP_IF_FALSE", HasArg: true},
	Break:       {Name: "BREAK"},
	Next:        {Name: "NEXT"},
	LoopStart:   {Name: "START_LOOP"},
	LoopEnd:     {Name: "END_LOOP"},

	MakeArray: {Name: "MAKE_ARRAY", HasArg: true},
	MakeTuple: {Name: "MAKE_TUPLE", HasArg: true},
	MakeMap:   {Name: "MAKE_MAP", HasArg: true},

	Use: {Name: "USE", HasArg: true},
}
