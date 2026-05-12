package main

import (
    "fmt"
)


type TokenType int

type Token struct {
	Kind  TokenType
	Value float64
}

const (
	Plus = iota
	Minus
	Multiply
	Caret
	IntNumber
	FloatNumber
	Variable
	VariableNeg
	Empty
    Equal
)

var TokenState = map[string]TokenType{
	"+":    Plus,
	"-":    Minus,
	"*":    Multiply,
	"^":    Caret,
	"X":    Variable,
	"(-X)": VariableNeg,
    "=": Equal,
}

var TokenStr = map[TokenType]string{
	Plus:        "+",
	Minus:       "-",
	Multiply:    "*",
	Caret:       "^",
	Variable:    "X",
	VariableNeg: "(-x)",
	Equal : "=",
}

func IsVariable(token Token) bool {

	t := token.GetType()
	if t != Variable && t != VariableNeg {
		return false
	}

	return true
}

type EquationToken interface {
	GetType() TokenType
	NumericValue() (float64, bool) // Returns value and a 'found' boolean
}

func (t *Token) GetType() TokenType {
	return t.Kind
}

func (t *Token) NumericValue() (float64, bool) {
	if t.Kind == IntNumber || t.Kind == FloatNumber {
		return t.Value, true
	}
	return 0, false
}

func (t Token) String() (string) {
	if t.Kind != FloatNumber {
		return fmt.Sprintf("%s", TokenStr[t.Kind]) 
	} 
	
    if t.Value < 0  {
        return fmt.Sprintf("(%.2f)", t.Value) 
    }
    return fmt.Sprintf("%.2f", t.Value) 
}

func (t *Token) ChangeSign() {
	if t.Kind == FloatNumber {
		t.Value = -t.Value
	} else if t.Kind == Variable {
		t.Kind = VariableNeg
	} else if t.Kind == VariableNeg {
		t.Kind = Variable
	}
}
