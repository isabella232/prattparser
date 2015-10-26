package main

import (
	"testing"
)

var testIntegerResult map[string]int = map[string]int{"1+4": 5,
	"14+4+5":          23,
	"4+5-7+3":         5,
	"5 + 3 * 4 - 2":   15,
	"5 + 3 * (4 - 2)": 11}

// Tests the evaluation of an expression that results in only some integer
func TestParsingForIntegerResult(t *testing.T) {
	for k, v := range testIntegerResult {
		if result, err := evalInt(k, nil); err == nil {
			if result.(int) != v {
				t.Errorf("Expected %d, actual %d for input %s", v, result.(int), k)
			}
		} else {
			t.Error("Unexpected evaluation error: ", err)
		}
	}
}

func evalInt(input string, variables map[string]int) (result interface{}, err error) {
	lexer := NewSimpleLexer(input)
	parser = &prattParser{lexer: lexer, variables: variables, expressionParsers: expressionParsers}
	parser.nextTokenParser, err = parser.nextTokenParserHandler()
	node, err := parser.parse(0)

	if err == nil {
		result = evaluateTree(node)
	}
	return
}

var testBoolResult map[string]bool = map[string]bool{"1 > 4": false,
	"14+4 > 5":                       true,
	"3-5 + 4 >= 2":                   true,
	"(4 > 3)":                        true,
	"(4 > 3) AND 3 <= 2 ":            false,
	"(4 > 3) AND 3 <= 2 OR  67 > 45": true,
	"4 > 6  OR NOT 4 > 8":            true}

func TestParsingForBooleanResult(t *testing.T) {
	for k, v := range testBoolResult {
		if result, err := eval(k, nil); err == nil {
			if result.(bool) != v {
				t.Errorf("Expected %d, actual %d  for input %s", v, result.(bool), k)
			}
		} else {
			t.Error("Unexpected evaluation error: ", err)
		}
	}
}

func TestParsingError(t *testing.T) {
	input := "8 > & 7"
	if _, err := eval(input, nil); err == nil {
		t.Errorf("Error expected, but not found")
	}

	input = " 6+3) > 7"
	if _, err := eval(input, nil); err == nil {
		t.Errorf("Error expected, but not found")
	}
}
