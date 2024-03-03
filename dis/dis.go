// Package dis supports analysis of Risor bytecode by disassembling it.
// This works with the opcodes defined in the `op` package and uses the
// InstructionIter type from the `compiler` package.
package dis

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/op"
)

var (
	bold      = color.New(color.Bold)
	yellow    = color.New(color.FgYellow)
	green     = color.New(color.FgGreen)
	magenta   = color.New(color.FgMagenta)
	italic    = color.New(color.Italic)
	nameColor = color.New(color.FgHiCyan)
)

// Instruction represents a single bytecode instruction and its operands.
type Instruction struct {
	Offset     int
	Name       string
	Opcode     op.Code
	Operands   []op.Code
	Annotation string
	Constant   interface{}
}

// Disassemble returns a parsed representation of the given bytecode.
func Disassemble(code *compiler.Code) ([]Instruction, error) {
	var instructions []Instruction
	var offset int
	iter := compiler.NewInstructionIter(code)
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		var err error
		info := op.GetInfo(val[0])
		var constant interface{}
		var annotation string
		switch info.Name {
		case "LOAD_FAST", "STORE_FAST":
			annotation, err = getLocalVariableName(code, int(val[1]))
			if err != nil {
				return nil, err
			}
		case "LOAD_GLOBAL", "STORE_GLOBAL":
			annotation, err = getGlobalVariableName(code, int(val[1]))
			if err != nil {
				return nil, err
			}
		case "LOAD_ATTR", "STORE_ATTR":
			nameIndex := int(val[1])
			name, err := getName(code, nameIndex)
			if err != nil {
				return nil, err
			}
			annotation = fmt.Sprintf("%v", name)
		case "BINARY_OP":
			annotation = op.BinaryOpType(val[1]).String()
		case "COMPARE_OP":
			annotation = op.CompareOpType(val[1]).String()
		case "LOAD_CONST":
			constant, err = getConstantValue(code, int(val[1]))
			if err != nil {
				return nil, err
			}
			annotation = fmt.Sprintf("%v", constant)
		}
		instructions = append(instructions, Instruction{
			Offset:     offset,
			Name:       info.Name,
			Opcode:     val[0],
			Operands:   val[1:],
			Annotation: annotation,
			Constant:   constant,
		})
		offset += len(val)
	}
	return instructions, nil
}

// Print a string representation of the given instructions to the given writer.
func Print(instructions []Instruction, writer io.Writer) {
	var lines [][]string
	for _, instr := range instructions {
		var values []string
		values = append(values, fmt.Sprintf("%d", instr.Offset))
		values = append(values, bold.Sprint(instr.Name))
		values = append(values, formatOperands(instr.Operands))
		if instr.Constant != nil {
			switch c := instr.Constant.(type) {
			case int64:
				values = append(values, yellow.Sprintf("%d", c))
			case float64:
				values = append(values, yellow.Sprintf("%f", c))
			case string:
				if len(c) > 80 {
					c = c[:77] + "..."
				}
				values = append(values, green.Sprintf("%q", c))
			case *compiler.Function:
				name := c.Name()
				if name == "" {
					name = italic.Sprint("<anonymous>")
				}
				values = append(values, magenta.Sprintf("func:%s", name))
			default:
				values = append(values, bold.Sprintf("%v", c))
			}
		} else if instr.Annotation != "" {
			values = append(values, nameColor.Sprintf("%v", instr.Annotation))
		} else {
			values = append(values, "")
		}
		lines = append(lines, values)
	}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"offset", "opcode", "operands", "info"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})
	for _, v := range lines {
		table.Append(v)
	}
	table.Render()
}

func formatOperands(ops []op.Code) string {
	var sb strings.Builder
	for i, op := range ops {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", op))
	}
	return sb.String()
}

func getLocalVariableName(code *compiler.Code, index int) (string, error) {
	if code.LocalsCount() <= index {
		return "", fmt.Errorf("local variable index out of range: %d", index)
	}
	return code.Local(index).Name(), nil
}

func getGlobalVariableName(code *compiler.Code, index int) (string, error) {
	if code.GlobalsCount() <= index {
		return "", fmt.Errorf("global variable index out of range: %d", index)
	}
	return code.Global(index).Name(), nil
}

func getConstantValue(code *compiler.Code, index int) (any, error) {
	if code.ConstantsCount() <= index {
		return "", fmt.Errorf("constant index out of range: %d", index)
	}
	return code.Constant(index), nil
}

func getName(code *compiler.Code, index int) (string, error) {
	if code.NameCount() <= index {
		return "", fmt.Errorf("name index out of range: %d", index)
	}
	return code.Name(index), nil
}
