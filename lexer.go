package main

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// simpleLexer is a simple lexer that holds the input string for tokenizing
type simpleLexer struct {
	input string
	start int
	pos   int
	width int
	err   error
}

const eof = -1

func NewSimpleLexer(src string) *simpleLexer {
	return &simpleLexer{input: src}
}

// next provides the next available token
func (lexer *simpleLexer) next() (tok *token, err error) {
	switch c := lexer.read(); {
	case c == eof:
		tok = &token{tokenTyp: END}
	case c == '(':
		tok = &token{tokenTyp: LEFTPAREN}
	case c == ')':
		tok = &token{tokenTyp: RIGHTPAREN}
	case c == '+':
		tok = &token{tokenTyp: PLUS}
	case c == '-':
		tok = &token{tokenTyp: MINUS}
	case c == '*':
		tok = &token{tokenTyp: MULTIPLY}
	case c == '>':
		if lexer.nextmatch('=') {
			tok = &token{tokenTyp: GREATEREQUAL}
		} else {
			tok = &token{tokenTyp: GREATER}
		}
	case c == '<':
		if lexer.nextmatch('=') {
			tok = &token{tokenTyp: LESSEREQUAL}
		} else {
			tok = &token{tokenTyp: LESSER}
		}
	case c == '=':
		if lexer.nextmatch('=') {
			tok = &token{tokenTyp: EQUAL}
		} else {
			tok = nil
		}
	case c == '!':
		if lexer.nextmatch('=') {
			tok = &token{tokenTyp: NOTEQUAL}
		} else {
			tok = nil
		}
	case isDigit(c):
		tok = lexer.extractNumber(c)
	case isIdentChar(c):
		tok = lexer.extractIdentifier(c)
	case isSpace(c):
		tok, err = lexer.skipSpace(c)
	default:
		//lexer.err = errors.New("Error: Unknown charater " + string(c))
		//fmt.Println("Unknown charater " + string(c))
		err = errors.New(fmt.Sprintf("Error: Invalid charater '%s' in the input\n", string(c)))
		tok = nil
	}
	return
}

func (lexer *simpleLexer) extractNumber(c rune) *token {
	b := make([]rune, 0)
	//extract number charaters till we encounter non numeric charater
	for ; isDigit(c); c = lexer.read() {
		b = append(b, c)
	}
	//Put back non number charater back
	lexer.backup()
	return &token{tokenTyp: NUMBER, value: string(b)}
}

func (lexer *simpleLexer) extractIdentifier(c rune) *token {
	b := make([]rune, 0)
	//extract alphanumeric characters till we encounter non alphanumeric charater
	for ; isIdentChar(c); c = lexer.read() {
		b = append(b, c)
	}
	//Put back non alphanumeric charater back
	lexer.backup()
	switch s := string(b); {
	case s == "AND":
		return &token{tokenTyp: AND}
	case s == "OR":
		return &token{tokenTyp: OR}
	case s == "NOT":
		return &token{tokenTyp: NOT}
	}
	return &token{tokenTyp: IDENTIFIER, value: string(b)}
}

func (lexer *simpleLexer) skipSpace(c rune) (*token, error) {
	for ; isSpace(c); c = lexer.read() {
	}
	//Put back non number charater back
	lexer.backup()
	return lexer.next()
}

// read reads one character at time and tracks the current position in the string
func (lexer *simpleLexer) read() rune {
	if lexer.pos >= len(lexer.input) {
		lexer.width = 0
		return eof
	}
	c, w := utf8.DecodeRuneInString(lexer.input[lexer.pos:])
	lexer.pos += w
	lexer.width = w
	return c
}

func (lexer *simpleLexer) backup() {
	lexer.pos -= lexer.width
}

func (lexer *simpleLexer) peek() rune {
	c := lexer.read()
	lexer.backup()
	return c
}

// nextmatch verifies whether the next character matches the given one
func (lexer *simpleLexer) nextmatch(c rune) bool {
	if c == lexer.read() {
		return true
	}
	lexer.backup()
	return false
}

func isIdentChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
