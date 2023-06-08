package op

type Code uint16

const (
	Nop Code = iota
	BinaryOp
	BinarySubscr
	BuildList
	BuildMap
	BuildSet
	BuildString
	Call
	CompareOp
	ContainsOp
	Copy
	DeleteSubscr
	False
	ForIter
	GetIter
	Halt
	Import
	JumpBackward
	JumpForward
	Length
	LoadAttr
	LoadBuiltin
	LoadClosure
	LoadConst
	LoadFast
	LoadFree
	LoadGlobal
	LoadName
	MakeCell
	Nil
	Partial
	PopJumpBackwardIfFalse
	PopJumpBackwardIfTrue
	PopJumpForwardIfFalse
	PopJumpForwardIfTrue
	PopTop
	Print
	PushNil
	Range
	ReturnValue
	Slice
	StoreAttr
	StoreFast
	StoreFree
	StoreGlobal
	StoreName
	StoreSubscr
	Swap
	True
	UnaryInvert
	UnaryNegative
	UnaryNot
	UnaryPositive
	Unpack
)

type BinaryOpType uint16

const (
	Add BinaryOpType = iota
	Subtract
	Multiply
	Divide
	Modulo
	And
	Or
	Xor
	Power
	LShift
	RShift
	BitwiseAnd
	BitwiseOr
)

type CompareOpType uint16

const (
	LessThan CompareOpType = iota
	LessThanOrEqual
	Equal
	NotEqual
	GreaterThan
	GreaterThanOrEqual
)

type Info struct {
	Name          string
	OperandCount  int
	OperandWidths []int
}

var OperandCount = make([]Info, 256)

func init() {
	type opInfo struct {
		op     Code
		name   string
		count  int
		widths []int
	}
	ops := []opInfo{
		{BinaryOp, "BINARY_OP", 1, []int{2}},
		{BinarySubscr, "BINARY_SUBSCR", 0, nil},
		{BuildList, "BUILD_LIST", 1, []int{2}},
		{BuildMap, "BUILD_MAP", 1, []int{2}},
		{BuildSet, "BUILD_SET", 1, []int{2}},
		{BuildString, "BUILD_STRING", 1, []int{2}},
		{Call, "CALL", 1, []int{2}},
		{CompareOp, "COMPARE_OP", 1, []int{2}},
		{ContainsOp, "CONTAINS_OP", 1, []int{2}},
		{Copy, "COPY", 1, []int{2}},
		{DeleteSubscr, "DELETE_SUBSCR", 0, nil},
		{False, "FALSE", 0, nil},
		{GetIter, "GET_ITER", 0, nil},
		{Halt, "HALT", 0, nil},
		{Import, "IMPORT", 0, nil},
		{JumpBackward, "JUMP_BACKWARD", 1, []int{2}},
		{JumpForward, "JUMP_FORWARD", 1, []int{2}},
		{Length, "LENGTH", 0, nil},
		{LoadAttr, "LOAD_ATTR", 1, []int{2}},
		{LoadBuiltin, "LOAD_BUILTIN", 1, []int{2}},
		{LoadClosure, "LOAD_CLOSURE", 2, []int{2, 2}},
		{LoadConst, "LOAD_CONST", 1, []int{2}},
		{LoadFast, "LOAD_FAST", 1, []int{2}},
		{LoadFree, "LOAD_FREE", 1, []int{2}},
		{LoadGlobal, "LOAD_GLOBAL", 1, []int{2}},
		{LoadName, "LOAD_NAME", 1, []int{2}},
		{MakeCell, "MAKE_CELL", 2, []int{2, 1}},
		{Nil, "NIL", 0, nil},
		{Nop, "NOP", 0, nil},
		{Partial, "PARTIAL", 1, []int{2}},
		{PopJumpBackwardIfFalse, "POP_JUMP_BACKWARD_IF_FALSE", 1, []int{2}},
		{PopJumpBackwardIfTrue, "POP_JUMP_BACKWARD_IF_TRUE", 1, []int{2}},
		{PopJumpForwardIfFalse, "POP_JUMP_FORWARD_IF_FALSE", 1, []int{2}},
		{PopJumpForwardIfTrue, "POP_JUMP_FORWARD_IF_TRUE", 1, []int{2}},
		{PopTop, "POP_TOP", 0, nil},
		{Print, "PRINT", 0, nil},
		{Range, "RANGE", 0, nil},
		{ReturnValue, "RETURN_VALUE", 0, nil},
		{Slice, "SLICE", 0, nil},
		{StoreAttr, "STORE_ATTR", 1, []int{2}},
		{StoreFast, "STORE_FAST", 1, []int{2}},
		{StoreFree, "STORE_FREE", 1, []int{2}},
		{StoreGlobal, "STORE_GLOBAL", 1, []int{2}},
		{StoreName, "STORE_NAME", 1, []int{2}},
		{StoreSubscr, "STORE_SUBSCR", 0, nil},
		{Swap, "SWAP", 1, []int{2}},
		{True, "TRUE", 0, nil},
		{UnaryNegative, "UNARY_NEGATIVE", 0, nil},
		{UnaryNot, "UNARY_NOT", 0, nil},
		{UnaryPositive, "UNARY_POSITIVE", 0, nil},
		{Unpack, "UNPACK", 1, []int{2}},
		{ForIter, "FOR_ITER", 2, []int{2, 2}},
	}
	for _, o := range ops {
		OperandCount[o.op] = Info{
			Name:          o.name,
			OperandCount:  o.count,
			OperandWidths: o.widths,
		}
	}
}

func GetInfo(op Code) Info {
	return OperandCount[op]
}
