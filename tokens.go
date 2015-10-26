package main

import (
	"fmt"
)

type TokenType int

// All token types that can be used in the input expressions
const (
	LEFTPAREN TokenType = iota
	RIGHTPAREN
	PLUS
	MINUS
	MULTIPLY
	GREATER
	GREATEREQUAL
	LESSER
	LESSEREQUAL
	NOTEQUAL
	EQUAL
	AND
	OR
	NOT
	NUMBER
	IDENTIFIER
	END
)

var tostring = map[TokenType]string{
	LEFTPAREN:    "LEFTPAREN",
	RIGHTPAREN:   "RIGHTPAREN",
	PLUS:         "PLUS",
	MINUS:        "MINUS",
	MULTIPLY:     "MULTIPLY",
	GREATER:      "GREATER",
	GREATEREQUAL: "GREATEREQUAL",
	LESSER:       "LESSER",
	LESSEREQUAL:  "LESSEREQUAL",
	NOTEQUAL:     "NOTEQUAL",
	EQUAL:        "EQUAL",
	AND:          "AND",
	OR:           "OR",
	NOT:          "NOT",
	NUMBER:       "NUMBER",
	IDENTIFIER:   "IDENTIFIER",
	END:          "END",
}

// tokenPrecedences maps a token to its precedence or binding power
var tokenPrecedences = map[TokenType]int{
	LEFTPAREN:    1,
	RIGHTPAREN:   1,
	PLUS:         40,
	MINUS:        40,
	MULTIPLY:     50,
	NUMBER:       0,
	GREATER:      30,
	GREATEREQUAL: 30,
	LESSER:       30,
	LESSEREQUAL:  30,
	NOTEQUAL:     30,
	EQUAL:        30,
	AND:          25,
	OR:           25,
	NOT:          0,
	IDENTIFIER:   0,
	END:          0,
}

// token represents a token with its type and value
type token struct {
	tokenTyp TokenType
	value    string
}

func (tok token) String() string {
	return fmt.Sprintf("tokenType:%s, value:%s", tostring[tok.tokenTyp], tok.value)
}
