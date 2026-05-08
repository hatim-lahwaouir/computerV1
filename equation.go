package main


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
	VariableDegree1
	VariableDegree2
	Empty
    VariableDegree2Negatif
    VariableDegree1Negatif
)

var TokenState = map[string]TokenType{
	"+":   Plus,
	"-":   Minus,
	"*":   Multiply,
	"^":   Caret,
	"X":   Variable,
	"x":   Variable,
	"x^2": VariableDegree2,
	"x^1": VariableDegree1,
	"(-x^2)": VariableDegree2Negatif,
	"(-x^1)": VariableDegree1Negatif,
}

var TokenStr = map[TokenType]string{
	Plus:            "+",
	Minus:           "-",
	Multiply:        "*",
	Caret:           "^",
	Variable:        "X",
	VariableDegree2: "x^2",
	VariableDegree1: "x^1",
    VariableDegree2Negatif: "(-x^2)",
    VariableDegree1Negatif: "(-x^1)",
}


func IsVariable(token Token) bool {

    t := token.GetType()
    if t != Variable && t != VariableDegree1 && t != VariableDegree2 && VariableDegree2Negatif != t && t != VariableDegree1Negatif {
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



func (t *Token) ChangeSign(){
    if t.Kind == FloatNumber {
        t.Value =  -t.Value
     }
}




