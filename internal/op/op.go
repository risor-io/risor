package op

type Code uint16

const (
	Nop Code = iota
	BinaryOp
	BinarySubscr
	BuildList
	BuildMap
	BuildSet
	Call
	CompareOp
	ContainsOp
	DeleteSubscr
	False
	Halt
	JumpBackward
	JumpForward
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
	PopJumpBackwardIfFalse
	PopJumpBackwardIfTrue
	PopJumpForwardIfFalse
	PopJumpForwardIfTrue
	PopTop
	Print
	PushNil
	ReturnValue
	StoreAttr
	StoreFast
	StoreFree
	StoreGlobal
	StoreName
	StoreSubscr
	True
	UnaryInvert
	UnaryNegative
	UnaryNot
	UnaryPositive
	Swap
	BuildString
	Partial
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
		{Call, "CALL", 1, []int{2}},
		{CompareOp, "COMPARE_OP", 1, []int{2}},
		{False, "FALSE", 0, nil},
		{Halt, "HALT", 0, nil},
		{PopJumpForwardIfTrue, "POP_JUMP_FORWARD_IF_TRUE", 1, []int{2}},
		{PopJumpBackwardIfTrue, "POP_JUMP_BACKWARD_IF_TRUE", 1, []int{2}},
		{PopJumpForwardIfFalse, "POP_JUMP_FORWARD_IF_FALSE", 1, []int{2}},
		{PopJumpBackwardIfFalse, "POP_JUMP_BACKWARD_IF_FALSE", 1, []int{2}},
		{JumpForward, "JUMP_FORWARD", 1, []int{2}},
		{JumpBackward, "JUMP_BACKWARD", 1, []int{2}},
		{LoadAttr, "LOAD_ATTR", 1, []int{2}},
		{LoadBuiltin, "LOAD_BUILTIN", 1, []int{2}},
		{LoadConst, "LOAD_CONST", 1, []int{2}},
		{LoadClosure, "LOAD_CLOSURE", 2, []int{2, 2}},
		{LoadFast, "LOAD_FAST", 1, []int{2}},
		{LoadGlobal, "LOAD_GLOBAL", 1, []int{2}},
		{LoadName, "LOAD_NAME", 1, []int{2}},
		{LoadFree, "LOAD_FREE", 1, []int{2}},
		{MakeCell, "MAKE_CELL", 2, []int{2, 1}},
		{Nil, "NIL", 0, nil},
		{Nop, "NOP", 0, nil},
		{PopTop, "POP_TOP", 0, nil},
		{Print, "PRINT", 0, nil},
		{ReturnValue, "RETURN_VALUE", 1, []int{2}},
		{StoreAttr, "STORE_ATTR", 1, []int{2}},
		{StoreFast, "STORE_FAST", 1, []int{2}},
		{StoreGlobal, "STORE_GLOBAL", 1, []int{2}},
		{StoreName, "STORE_NAME", 1, []int{2}},
		{StoreFree, "STORE_FREE", 1, []int{2}},
		{True, "TRUE", 0, nil},
		{UnaryNegative, "UNARY_NEGATIVE", 0, nil},
		{UnaryNot, "UNARY_NOT", 0, nil},
		{UnaryPositive, "UNARY_POSITIVE", 0, nil},
		{BuildList, "BUILD_LIST", 1, []int{2}},
		{BuildMap, "BUILD_MAP", 1, []int{2}},
		{BuildSet, "BUILD_SET", 1, []int{2}},
		{BinarySubscr, "BINARY_SUBSCR", 0, nil},
		{DeleteSubscr, "DELETE_SUBSCR", 0, nil},
		{StoreSubscr, "STORE_SUBSCR", 0, nil},
		{ContainsOp, "CONTAINS_OP", 1, []int{2}},
		{Swap, "SWAP", 1, []int{2}},
		{BuildString, "BUILD_STRING", 1, []int{2}},
		{Partial, "PARTIAL", 1, []int{2}},
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

// Python/ceval.c
// https://bytecode.readthedocs.io/en/latest/api.html#binary-operation
// static const binaryfunc binary_ops[] = {
//     [NB_ADD] = PyNumber_Add,
//     [NB_AND] = PyNumber_And,
//     [NB_FLOOR_DIVIDE] = PyNumber_FloorDivide,
//     [NB_LSHIFT] = PyNumber_Lshift,
//     [NB_MATRIX_MULTIPLY] = PyNumber_MatrixMultiply,
//     [NB_MULTIPLY] = PyNumber_Multiply,
//     [NB_REMAINDER] = PyNumber_Remainder,
//     [NB_OR] = PyNumber_Or,
//     [NB_POWER] = _PyNumber_PowerNoMod,
//     [NB_RSHIFT] = PyNumber_Rshift,
//     [NB_SUBTRACT] = PyNumber_Subtract,
//     [NB_TRUE_DIVIDE] = PyNumber_TrueDivide,
//     [NB_XOR] = PyNumber_Xor,
//     [NB_INPLACE_ADD] = PyNumber_InPlaceAdd,
//     [NB_INPLACE_AND] = PyNumber_InPlaceAnd,
//     [NB_INPLACE_FLOOR_DIVIDE] = PyNumber_InPlaceFloorDivide,
//     [NB_INPLACE_LSHIFT] = PyNumber_InPlaceLshift,
//     [NB_INPLACE_MATRIX_MULTIPLY] = PyNumber_InPlaceMatrixMultiply,
//     [NB_INPLACE_MULTIPLY] = PyNumber_InPlaceMultiply,
//     [NB_INPLACE_REMAINDER] = PyNumber_InPlaceRemainder,
//     [NB_INPLACE_OR] = PyNumber_InPlaceOr,
//     [NB_INPLACE_POWER] = _PyNumber_InPlacePowerNoMod,
//     [NB_INPLACE_RSHIFT] = PyNumber_InPlaceRshift,
//     [NB_INPLACE_SUBTRACT] = PyNumber_InPlaceSubtract,
//     [NB_INPLACE_TRUE_DIVIDE] = PyNumber_InPlaceTrueDivide,
//     [NB_INPLACE_XOR] = PyNumber_InPlaceXor,
// };
