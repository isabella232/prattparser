package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// prattParser represents the main parser
type prattParser struct {
	lexer             *simpleLexer
	nextTokenParser   tokenParser
	expressionParsers map[TokenType]expressionParser
	variables         map[string]int
}

// tokenParser has a token and its expressionParser where
// expressionParser is the interface that wraps the token parsing methods
type tokenParser struct {
	expressionParser
	token token
}

// parse parses the scanned tokens
func (pp *prattParser) parse(startPrecedence int) (left *parseNode, err error) {
	var current = pp.nextTokenParser
	pp.nextTokenParser, err = pp.nextTokenParserHandler()
	if err != nil {
		return
	}
	left, err = current.parseExpression(pp, current.token)
	if err != nil {
		return
	}
	for startPrecedence < tokenPrecedences[pp.nextTokenParser.token.tokenTyp] {
		current = pp.nextTokenParser
		pp.nextTokenParser, err = pp.nextTokenParserHandler()
		if err != nil {
			return
		}
		left, err = current.parseExpressionWithLeft(pp, current.token, left)
		if err != nil {
			return
		}
	}
	return
}

// nextTokenParserHandler gets the next token from lexer and returns "tokenParser"
// which has the next token and the token's "expressionParser"
func (pp *prattParser) nextTokenParserHandler() (tokParser tokenParser, err error) {
	var tok *token
	if tok, err = pp.lexer.next(); err == nil {
		// gets the token's expressionParser from the mapping
		if handler, ok := pp.expressionParsers[tok.tokenTyp]; ok {
			tokParser = tokenParser{handler, *tok}
		} else {
			tokParser = tokenParser{dummyExpressionParser{}, *tok}
		}
	} else {
		tokParser = tokenParser{dummyExpressionParser{}, token{tokenTyp: END}}
	}
	return
}

func (pp *prattParser) matchAndAdvance(tokenTyp TokenType) (err error) {
	if pp.nextTokenParser.token.tokenTyp != tokenTyp {
		return fmt.Errorf("Error: Expected %s, but found  %s", tostring[tokenTyp], tostring[pp.nextTokenParser.token.tokenTyp])
	}
	pp.nextTokenParser, err = pp.nextTokenParserHandler()
	return
}

var parser *prattParser = nil

func eval(input string, variables map[string]int) (result interface{}, err error) {

	lexer := NewSimpleLexer(input)
	parser = &prattParser{lexer: lexer, variables: variables, expressionParsers: expressionParsers}
	parser.nextTokenParser, err = parser.nextTokenParserHandler()
	node, err := parser.parse(0)

	if err == nil {
		v := evaluateTree(node)
		switch v.(type) {
		case bool:
			result = v
		case int:
			iv := v.(int)
			result = iv > 0
		}
	}
	return
}

var (
	inputExpression = flag.String("input", "", "An expression string which we want to evaluate")
	inputVariables  = flag.String("variables", "", "List of key value pairs. Eg VAR1=100,VAR2=45 ")
)

func extractVariables(keyvalStr string) map[string]int {
	varibales := make(map[string]int, 0)
	if len(strings.Trim(keyvalStr, " ")) == 0 {
		return varibales
	}
	for _, v := range strings.Split(keyvalStr, ",") {
		kv := strings.Split(strings.Trim(v, " "), "=")
		if vv, err := strconv.Atoi(kv[1]); err == nil {
			varibales[kv[0]] = vv
		} else {
			log.Fatalf(" Integer conversion error ", err)
		}
	}
	return varibales
}

func main() {
	flag.Parse()
	if len(*inputExpression) == 0 {
		log.Fatalf("Missing required --input parameter")
	}

	varibales := extractVariables(*inputVariables)

	if v, err := eval(*inputExpression, varibales); err == nil {
		fmt.Println("evaluated result is ", v)
	} else {
		fmt.Printf("%s\n", err)
	}
}
