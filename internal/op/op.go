package op

type Code byte

const (
	Nop                    Code = 1
	LoadConst              Code = 2
	LoadName               Code = 3
	LoadGlobal             Code = 4
	LoadFast               Code = 5 // Pushes a reference to the local co_varnames[var_num] onto the stack.
	LoadAttr               Code = 6
	StoreFast              Code = 7 // Stores TOS into the local co_varnames[var_num].
	StoreName              Code = 8
	StoreAttr              Code = 9
	StoreGlobal            Code = 10
	Nil                    Code = 11
	True                   Code = 12
	False                  Code = 13
	BinaryOp               Code = 14
	UnaryPositive          Code = 15
	UnaryNegative          Code = 16
	UnaryNot               Code = 17
	UnaryInvert            Code = 18
	CompareOp              Code = 19
	ReturnValue            Code = 20
	PopJumpForwardIfTrue   Code = 21
	PopJumpBackwardIfTrue  Code = 22
	PopJumpForwardIfFalse  Code = 23
	PopJumpBackwardIfFalse Code = 24
	Halt                   Code = 25
	Print                  Code = 26
	PopTop                 Code = 27
	PushNil                Code = 28
	Call                   Code = 29
	JumpForward            Code = 30
	JumpBackward           Code = 31
	BuildList              Code = 32
	BuildMap               Code = 33
)

type BinaryOpType byte

const (
	Add      BinaryOpType = 1
	Subtract BinaryOpType = 2
	Multiply BinaryOpType = 3
	Divide   BinaryOpType = 4
	Modulo   BinaryOpType = 5
	And      BinaryOpType = 6
	Or       BinaryOpType = 7
	Xor      BinaryOpType = 8
	Power    BinaryOpType = 9
	LShift   BinaryOpType = 10
	RShift   BinaryOpType = 11
)

type CompareOpType byte

const (
	LessThan           CompareOpType = 1
	LessThanOrEqual    CompareOpType = 2
	Equal              CompareOpType = 3
	NotEqual           CompareOpType = 4
	GreaterThan        CompareOpType = 5
	GreaterThanOrEqual CompareOpType = 6
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
		{BinaryOp, "BINARY_OP", 1, []int{1}},
		{Call, "CALL", 1, []int{1}},
		{CompareOp, "COMPARE_OP", 1, []int{1}},
		{False, "FALSE", 0, nil},
		{Halt, "HALT", 0, nil},
		{PopJumpForwardIfTrue, "POP_JUMP_FORWARD_IF_TRUE", 1, []int{2}},
		{PopJumpBackwardIfTrue, "POP_JUMP_BACKWARD_IF_TRUE", 1, []int{2}},
		{PopJumpForwardIfFalse, "POP_JUMP_FORWARD_IF_FALSE", 1, []int{2}},
		{PopJumpBackwardIfFalse, "POP_JUMP_BACKWARD_IF_FALSE", 1, []int{2}},
		{JumpForward, "JUMP_FORWARD", 1, []int{2}},
		{JumpBackward, "JUMP_BACKWARD", 1, []int{2}},
		{LoadAttr, "LOAD_ATTR", 1, []int{2}},
		{LoadConst, "LOAD_CONST", 1, []int{2}},
		{LoadFast, "LOAD_FAST", 1, []int{2}},
		{LoadGlobal, "LOAD_GLOBAL", 1, []int{2}},
		{LoadName, "LOAD_NAME", 1, []int{2}},
		{Nil, "NIL", 0, nil},
		{Nop, "NOP", 0, nil},
		{PopTop, "POP_TOP", 0, nil},
		{Print, "PRINT", 0, nil},
		{ReturnValue, "RETURN_VALUE", 1, []int{1}},
		{StoreAttr, "STORE_ATTR", 1, []int{1}},
		{StoreFast, "STORE_FAST", 1, []int{1}},
		{StoreGlobal, "STORE_GLOBAL", 1, []int{2}},
		{StoreName, "STORE_NAME", 1, []int{1}},
		{True, "TRUE", 0, nil},
		{UnaryNegative, "UNARY_NEGATIVE", 0, nil},
		{UnaryNot, "UNARY_NOT", 0, nil},
		{UnaryPositive, "UNARY_POSITIVE", 0, nil},
		{UnaryPositive, "UNARY_POSITIVE", 0, nil},
		{BuildList, "BUILD_LIST", 1, []int{2}},
		{BuildMap, "BUILD_MAP", 1, []int{2}},
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
