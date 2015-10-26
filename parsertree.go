package main

import (
	"fmt"
	"strconv"
)

// parseNode represents each parsed token expression
type parseNode struct {
	token token
	left  *parseNode
	right *parseNode
}

// expressionParser is the interface that wraps the token parsing methods
type expressionParser interface {
	parseExpression(parser *prattParser, tok token) (*parseNode, error)
	parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error)
}

// dummyExpressionParser is an expressionParser implementation that we use whenever no token parser is available in the mapping
type dummyExpressionParser struct{}

func (p dummyExpressionParser) parseExpression(parser *prattParser, tok token) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpression")
}

func (p dummyExpressionParser) parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpressionWithLeft")
}

// defaultPrefixParser implements only the parseExpression method
type defaultPrefixParser struct{}

func (p defaultPrefixParser) parseExpression(parser *prattParser, tok token) (*parseNode, error) {
	return &parseNode{tok, nil, nil}, nil
}

func (p defaultPrefixParser) parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpressionWithLeft")
}

// nodeInfixExpressionParser implements only the parseExpressionWithLeft method
type nodeInfixExpressionParser struct{}

func (p nodeInfixExpressionParser) parseExpression(parser *prattParser, tok token) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpression")
}

func (p nodeInfixExpressionParser) parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error) {
	right, err := parser.parse(tokenPrecedences[tok.tokenTyp])
	return &parseNode{tok, left, right}, err
}

type nodePrefixExpressionParser struct{}

func (p nodePrefixExpressionParser) parseExpression(parser *prattParser, tok token) (*parseNode, error) {
	left, err := parser.parse(tokenPrecedences[tok.tokenTyp])
	return &parseNode{tok, left, nil}, err
}

func (p nodePrefixExpressionParser) parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpressionWithLeft")
}

// leftparenNodePrefixParser implements only the parseExpression method as leftparen will need to concern about all tokens coming after it
type leftparenNodePrefixParser struct{}

func (p leftparenNodePrefixParser) parseExpression(parser *prattParser, tok token) (*parseNode, error) {
	left, err := parser.parse(tokenPrecedences[tok.tokenTyp])
	if err == nil {
		err = parser.matchAndAdvance(RIGHTPAREN)
	}
	//return &parseNode{tok, left, nil}
	return left, err
}

func (p leftparenNodePrefixParser) parseExpressionWithLeft(parser *prattParser, tok token, left *parseNode) (*parseNode, error) {
	return nil, fmt.Errorf("Token %s does not support method '%s'", tostring[tok.tokenTyp], "parseExpressionWithLeft")
}

// expressionParsers has mapping of tokens to the expressionParser
var expressionParsers = map[TokenType]expressionParser{
	LEFTPAREN:    leftparenNodePrefixParser{},
	PLUS:         nodeInfixExpressionParser{},
	MINUS:        nodeInfixExpressionParser{},
	IDENTIFIER:   defaultPrefixParser{},
	NUMBER:       defaultPrefixParser{},
	MULTIPLY:     nodeInfixExpressionParser{},
	EQUAL:        nodeInfixExpressionParser{},
	NOTEQUAL:     nodeInfixExpressionParser{},
	GREATER:      nodeInfixExpressionParser{},
	GREATEREQUAL: nodeInfixExpressionParser{},
	LESSER:       nodeInfixExpressionParser{},
	LESSEREQUAL:  nodeInfixExpressionParser{},
	AND:          nodeInfixExpressionParser{},
	OR:           nodeInfixExpressionParser{},
	NOT:          nodePrefixExpressionParser{},
}

func evaluateTree(n *parseNode) (result interface{}) {
	switch n.token.tokenTyp {
	case NUMBER:
		result, _ = strconv.Atoi(n.token.value)
	case PLUS:
		result = evaluateTree(n.left).(int) + evaluateTree(n.right).(int)
	case MINUS:
		result = evaluateTree(n.left).(int) - evaluateTree(n.right).(int)
	case MULTIPLY:
		result = evaluateTree(n.left).(int) * evaluateTree(n.right).(int)
	case IDENTIFIER:
		result = parser.variables[n.token.value]
	case EQUAL:
		result = evaluateTree(n.left).(int) == evaluateTree(n.right).(int)
	case NOTEQUAL:
		result = evaluateTree(n.left).(int) != evaluateTree(n.right).(int)
	case GREATER:
		result = evaluateTree(n.left).(int) > evaluateTree(n.right).(int)
	case GREATEREQUAL:
		result = evaluateTree(n.left).(int) >= evaluateTree(n.right).(int)
	case LESSER:
		result = evaluateTree(n.left).(int) < evaluateTree(n.right).(int)
	case LESSEREQUAL:
		result = evaluateTree(n.left).(int) <= evaluateTree(n.right).(int)
	case AND:
		result = evaluateTree(n.left).(bool) && evaluateTree(n.right).(bool)
	case OR:
		result = evaluateTree(n.left).(bool) || evaluateTree(n.right).(bool)
	case NOT:
		result = !evaluateTree(n.left).(bool)
	default:
		result = nil
	}
	return
}
