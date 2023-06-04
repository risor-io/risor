package parser

import "github.com/cloudcmds/tamarin/v2/token"

// Precedence order for operators
const (
	_ int = iota
	LOWEST
	PIPE        // |
	COND        // OR or AND
	ASSIGN      // =
	DECLARE     // :=
	TERNARY     // ? :
	EQUALS      // == or !=
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	POWER       // **
	MOD         // %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	IN          // X in Y
	RANGE       // range X
	INDEX       // array[index], map[key]
	HIGHEST
)

// Precedences for each token type
var precedences = map[token.Type]int{
	token.QUESTION:        TERNARY,
	token.ASSIGN:          ASSIGN,
	token.DECLARE:         DECLARE,
	token.EQ:              EQUALS,
	token.NOT_EQ:          EQUALS,
	token.LT:              LESSGREATER,
	token.LT_EQUALS:       LESSGREATER,
	token.GT:              LESSGREATER,
	token.GT_EQUALS:       LESSGREATER,
	token.PLUS:            SUM,
	token.PLUS_EQUALS:     SUM,
	token.MINUS:           SUM,
	token.MINUS_EQUALS:    SUM,
	token.SLASH:           PRODUCT,
	token.SLASH_EQUALS:    PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.ASTERISK_EQUALS: PRODUCT,
	token.POW:             POWER,
	token.MOD:             MOD,
	token.AND:             COND,
	token.OR:              COND,
	token.PIPE:            PIPE,
	token.LPAREN:          CALL,
	token.PERIOD:          CALL,
	token.LBRACKET:        INDEX,
	token.IN:              IN,
	token.RANGE:           RANGE,
}
