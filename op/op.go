// Package op defines the opcodes that are used in the Risor virtual machine.
package op

// Code is an integer opcode that indicates an operation to execute.
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

// BinaryOpType describes a type of binary operation.
type BinaryOpType uint16

const (
	Add BinaryOpType = iota + 1
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

// CompareOpType describes a type of comparison operation.
type CompareOpType uint16

const (
	LessThan CompareOpType = iota + 1
	LessThanOrEqual
	Equal
	NotEqual
	GreaterThan
	GreaterThanOrEqual
)

// Info contains information about an opcode.
type Info struct {
	Name         string
	OperandCount int
}

var infos = make([]Info, 256)

func init() {
	type opInfo struct {
		op    Code
		name  string
		count int
	}
	ops := []opInfo{
		{BinaryOp, "BINARY_OP", 1},
		{BinarySubscr, "BINARY_SUBSCR", 0},
		{BuildList, "BUILD_LIST", 1},
		{BuildMap, "BUILD_MAP", 1},
		{BuildSet, "BUILD_SET", 1},
		{BuildString, "BUILD_STRING", 1},
		{Call, "CALL", 1},
		{CompareOp, "COMPARE_OP", 1},
		{ContainsOp, "CONTAINS_OP", 1},
		{Copy, "COPY", 1},
		{DeleteSubscr, "DELETE_SUBSCR", 0},
		{False, "FALSE", 0},
		{GetIter, "GET_ITER", 0},
		{Halt, "HALT", 0},
		{Import, "IMPORT", 0},
		{JumpBackward, "JUMP_BACKWARD", 1},
		{JumpForward, "JUMP_FORWARD", 1},
		{Length, "LENGTH", 0},
		{LoadAttr, "LOAD_ATTR", 1},
		{LoadClosure, "LOAD_CLOSURE", 2},
		{LoadConst, "LOAD_CONST", 1},
		{LoadFast, "LOAD_FAST", 1},
		{LoadFree, "LOAD_FREE", 1},
		{LoadGlobal, "LOAD_GLOBAL", 1},
		{LoadName, "LOAD_NAME", 1},
		{MakeCell, "MAKE_CELL", 2},
		{Nil, "NIL", 0},
		{Nop, "NOP", 0},
		{Partial, "PARTIAL", 1},
		{PopJumpBackwardIfFalse, "POP_JUMP_BACKWARD_IF_FALSE", 1},
		{PopJumpBackwardIfTrue, "POP_JUMP_BACKWARD_IF_TRUE", 1},
		{PopJumpForwardIfFalse, "POP_JUMP_FORWARD_IF_FALSE", 1},
		{PopJumpForwardIfTrue, "POP_JUMP_FORWARD_IF_TRUE", 1},
		{PopTop, "POP_TOP", 0},
		{Print, "PRINT", 0},
		{Range, "RANGE", 0},
		{ReturnValue, "RETURN_VALUE", 0},
		{Slice, "SLICE", 0},
		{StoreAttr, "STORE_ATTR", 1},
		{StoreFast, "STORE_FAST", 1},
		{StoreFree, "STORE_FREE", 1},
		{StoreGlobal, "STORE_GLOBAL", 1},
		{StoreName, "STORE_NAME", 1},
		{StoreSubscr, "STORE_SUBSCR", 0},
		{Swap, "SWAP", 1},
		{True, "TRUE", 0},
		{UnaryNegative, "UNARY_NEGATIVE", 0},
		{UnaryNot, "UNARY_NOT", 0},
		{UnaryPositive, "UNARY_POSITIVE", 0},
		{Unpack, "UNPACK", 1},
		{ForIter, "FOR_ITER", 2},
	}
	for _, o := range ops {
		infos[o.op] = Info{
			Name:         o.name,
			OperandCount: o.count,
		}
	}
}

// GetInfo returns information about the given opcode.
func GetInfo(op Code) Info {
	return infos[op]
}
