package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)




// error function

func Fatal(msg string) {

	fmt.Printf("error: %s\n", msg)

	os.Exit(1)
}

// helpers

func IsOperation(token Token) bool {
    t := token.GetType()
	if t != Plus && t != Minus && t != Multiply && Caret != t {
		return false
	}

	return true
}




func SimplifyVariables1(tokens []Token) ([]Token, error) {
	var (
		newTokens   []Token
        number float64
	)

	for i := 0; i < len(tokens); i++ {

		if tokens[i].GetType() == Variable {
			if i+1 == len(tokens) {
				newTokens = append(newTokens, Token{Kind: VariableDegree1})
			} else {

				if IsOperation(tokens[i+1]) == false {
					return nil, errors.New("Invalid equation operations must be after the variable")
				}

				if tokens[i+1].GetType() == Caret {
                    number, _ = tokens[i+2].NumericValue()
					switch  number {
					case 0:
						newTokens = append(newTokens, Token{Kind: FloatNumber, Value: 1})
					case 1:
				        newTokens = append(newTokens, Token{Kind: VariableDegree1})
					case 2:
				        newTokens = append(newTokens, Token{Kind: VariableDegree2})
					default:
						return nil, errors.New("we only solve second degree equations")
					}
					i += 2
				} else {
				    newTokens = append(newTokens, Token{Kind: VariableDegree1})
				}
			}

		} else {
			newTokens = append(newTokens, tokens[i])
		}
	}

	return newTokens,  nil
}

func SimplifyVariables2(tokens []Token) ([]Token, error) {

	var (
		newTokens []Token
	)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].GetType() == Empty {
			continue
		}

		if tokens[i].GetType() == VariableDegree1 || tokens[i].GetType() == VariableDegree2 {
			var (
				oper   bool
				degree int
			)
			oper = false
			degree = 0

			if i+1 >= len(tokens) {
				newTokens = append(newTokens, tokens[i])
				break
			}
			j := i
			for ; j < len(tokens); j++ {
				if oper {
					if tokens[j].GetType() != Multiply {
						break
					}
				} else {
					if degree == 1 && (tokens[j].GetType() == VariableDegree2 || tokens[j].GetType() == VariableDegree1) {
						tokens[j-1] = Token{Kind: Empty}
					}

					if tokens[j].GetType() == VariableDegree1 {
						tokens[j] = Token{Kind: Empty}
						degree++
					} else if tokens[j].GetType() == VariableDegree2 {
						tokens[j] = Token{Kind: Empty}
						degree += 2
					}

				}

				oper = !oper
			}

			if degree == 1 {
				newTokens = append(newTokens, Token{Kind: VariableDegree1})
			} else if degree == 2 {
				newTokens = append(newTokens, Token{Kind: VariableDegree2})
			} else {
				return nil, errors.New("we only solve second degree equations")
			}

		} else {
			newTokens = append(newTokens, tokens[i])
		}
	}
	return newTokens, nil

}

func SimplifyMultiplication(tokens []Token) ([]Token, error) {
	var (
		newTokens   []Token
        number float64
	)

	for i := 0; i < len(tokens); i++ {

		if tokens[i].GetType() == Empty {
			continue
		}

		if tokens[i].GetType() == FloatNumber {

			var (
				oper   bool
				result float64
			)
			oper = false
			result = 1

			for j := i; j < len(tokens); j++ {

				if oper {
					if tokens[j].GetType() != Multiply {
						break
					}
				}

				if tokens[j].GetType() == FloatNumber {
                  
			        number, _ = tokens[j]. NumericValue()
					result = result * number 
                    tokens[j] = Token{Kind: Empty}
					if j > 0 && tokens[j-1].GetType() == Multiply {
						tokens[j-1] = Token{Kind : Empty}
					}
				}

				oper = !oper
			}
			newTokens = append(newTokens, Token{Kind: FloatNumber, Value: result})

		} else {
			newTokens = append(newTokens, tokens[i])
		}
	}
	return newTokens, nil

}


func SimplifyAddition(tokens []Token) ([]Token, error) {

    
    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == FloatNumber {
            if i > 0 && tokens[i - 1].GetType() == Minus{
                tokens[i - 1] = Token{Kind : Plus} 
                tokens[i].ChangeSign() 
            }
        }
    }

    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == VariableDegree2 || tokens[i].GetType() == VariableDegree1 {

            if i > 0 && tokens[i - 1].GetType() == Minus{
                tokens[i - 1] = Token{Kind: Plus}
                if VariableDegree2 == tokens[i].GetType() {
                    tokens[i] = Token{Kind: VariableDegree2Negatif}

                }else {
                    tokens[i] = Token{Kind: VariableDegree1Negatif}
                }
            }
        }
    }

    return tokens, nil

}


func ReduceVersion(tokens []Token ) ([]Token, error) {
    var (
        equ map[string]float64
        state TokenType
        number float64
    )

    equ = make(map[string]float64)

    equ["X^1"] = 0
    equ["X^2"] = 0
    equ["C"] = 0


    for i := 0; i < len(tokens); i++{

        if tokens[i].GetType() != Multiply {
            continue
        }
        
        if IsVariable(tokens[i - 1]) {
            state = tokens[i -1].GetType()
            number, _ = tokens[i + 1].NumericValue() 
        }else{
            state = tokens[i + 1].GetType()
            number, _ = tokens[i - 1].NumericValue() 
        }

        tokens[i + 1] = Token{Kind: Empty}
        tokens[i - 1] = Token{Kind: Empty}
        tokens[i] = Token{Kind: Empty}

        if state == VariableDegree1 {
            equ["X^1"] += number 
        }
        if state == VariableDegree2 {
            equ["X^2"] += number 
        }
    
        if state == VariableDegree2Negatif {
            if number < 0 {
                equ["X^2"] += (-number)
            }else {
                equ["X^2"] += number
            }
        }

        if state == VariableDegree1Negatif {
            if number < 0 {
                equ["X^1"] += (-number)
            }else {
                equ["X^1"] += number
            }
        }
    }

    for i := 0; i < len(tokens); i++{
        if IsVariable(tokens[i]) == false{
            continue
        }

        switch tokens[i].GetType(){
            case VariableDegree1:
                equ["X^1"] += 1 
            case VariableDegree2:
                equ["X^2"] += 1 
            case VariableDegree2Negatif:
                equ["X^2"] -= 1 
            case VariableDegree1Negatif:
                equ["X^1"] -= 1 
         }

        tokens[i] = Token{Kind: Empty}
        }

    for i := 0; i < len(tokens); i++{
        if nbr, ok := tokens[i].NumericValue(); ok == true{
            equ["c"]+= nbr
        }
    }

    fmt.Printf("X^1 * %.2f + X ^2 * %.2f + (%.2f)", equ["X^1"],equ["X^2"], equ["c"])

    return nil, nil 
}



func Parse(input string) ([]Token, error) {
	var (
		tokens       []Token
	)

	for i := 0; i < len(input); i++ {

		if unicode.IsSpace(rune(input[i])) {
			continue
		}

		if _ , ok := TokenState[string(input[i])]; ok == true {
			tokens = append(tokens, Token{Kind: TokenState[string(input[i])]})

		} else {
			// treating the number
			var (
				s       int
				e       int
				isFloat bool
			)

			isFloat = false
			s = i
			for ; i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.'); i++ {
				if input[i] == '.' {
					isFloat = true

				}

			}
			e = i
			if isFloat == true {
				val, err := strconv.ParseFloat(input[s:e], 32)

				if err != nil {
					fmt.Println(input[s:e], s, e)
					Fatal("error parssing number")
				}
				tokens = append(tokens, Token{Kind: FloatNumber, Value: float64(val) })
			} else {
				val, err := strconv.Atoi(input[s:e])
				if err != nil {
					fmt.Println(input[s:e], s, e)
					Fatal("error parssing number")
				}
				tokens = append(tokens, Token{Kind: FloatNumber, Value: float64(val) })
			}
			i--
		}
	}

	return tokens, nil
}

func dbg(tokens []Token) {

	fmt.Printf("dbgV1:")
	for i := 0; i < len(tokens); i++ {
		val, ok := TokenStr[tokens[i].GetType()]
		if ok == true {
			fmt.Printf(" %s ", val)
		} else  {
            nbr, _ := tokens[i].NumericValue() 
			fmt.Printf(" %f ", nbr)
		}
	}
	fmt.Println()

}


func PrintTokens(tokens []Token) {

	fmt.Printf("dbgV2:")
	for i := 0; i < len(tokens); i++ {
		val, ok := TokenStr[tokens[i].GetType()]
		if ok == true {
			fmt.Printf(" %s ", val)
		} else {
            nbr, _ := tokens[i].NumericValue() 
            if nbr < 0 {
			    fmt.Printf(" (%.2f) ", nbr)
            }else {
			fmt.Printf(" %.2f ", nbr)
            }
		}
	}

	fmt.Println()

}

func dbgV2(tokens []Token) {


	fmt.Printf("dbgV2:")
	for i := 0; i < len(tokens); i++ {
		val, ok := TokenStr[tokens[i].GetType()]
		if ok == true {
			fmt.Printf(" %s ", val)
		} else {

            nbr, _ := tokens[i].NumericValue() 
            if nbr < 0 {
			    fmt.Printf(" (%.2f) ", nbr)
            }else {
			fmt.Printf(" %.2f ", nbr)
            }
		}
	}

	fmt.Println()

}

func main() {

	var (
		tokens       []Token
	)

	if len(os.Args) != 2 {
		Fatal("we need 2 arguments")
	}

	tokens, _ = Parse(os.Args[1])

	tokens,  err := SimplifyVariables1(tokens)

	if err != nil {
		Fatal(err.Error())
	}

	tokens, err = SimplifyVariables2(tokens)
	if err != nil {
		Fatal(err.Error())
	}
	dbg(tokens)

	tokens,  _ = SimplifyMultiplication(tokens )
    tokens,_ = SimplifyAddition(tokens)

	dbgV2(tokens)
    tokens, _ = ReduceVersion(tokens)

}
