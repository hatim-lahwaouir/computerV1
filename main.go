package main

import (
//	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
	"sort"
)






// F = X | Number | X * CARET  
func ParseFactor(token []Token, cur *int) (TreeNode){

    if *cur >= len(token) {
        return nil

    }

    var (
        t TokenType
    )

    // we still don't know if we need to validate negatif number
    t = token[*cur].GetType()

    if t !=  FloatNumber && t !=  Variable && t!= VariableNeg{
        return nil
    }


    
    if t ==  FloatNumber{
        val, _ := token[*cur].NumericValue()
        (*cur)++
		if (*cur) + 1 < len(token) && token[*cur].GetType() == Caret {
			(*cur)++
			return NewOpr(NewVal(val), ParseFactor(token, cur), Caret)
		}
        return NewVal(val)
    } 

	// check for ^ 
    if t ==  Variable {
        (*cur)++
	
		if (*cur) + 1 < len(token) && token[*cur].GetType() == Caret {
			(*cur)++
			return NewOpr(NewVar(1, 1), ParseFactor(token, cur), Caret)
		}
        return  NewVar(1, 1)
    }

    if t ==  VariableNeg {
        (*cur)++
		if (*cur) + 1 < len(token) && token[*cur].GetType() == Caret {
			(*cur)++
			return NewOpr(NewVar(1, -1), ParseFactor(token, cur), Caret)
		}
        return  NewVar(1, -1)
    }
    
    return nil
}

// P = F ^  NUMBER


// T = F * T |  F  
func ParseTerm(token []Token, cur *int) (TreeNode){

    var (
        a TreeNode
        b TreeNode
        t TokenType
    )

    a = ParseFactor(token, cur)

    if *cur >= len(token) {
        return a
    }

    t = token[*cur].GetType()

    if t == Multiply {
        (*cur)++
        b = ParseTerm(token, cur)

        return NewOpr(a, b, Multiply)
    } else if t == Caret{
        (*cur)++
        b = ParseTerm(token, cur )

        return NewOpr(a, b, Caret)
    }

    return a 
}


// E =  T {+|- T} 

func ParseExpression(token []Token, cur *int) (TreeNode){

    var (
        a TreeNode
        b TreeNode
        t TokenType
    )

    a = ParseTerm(token, cur)

    for ;; {
        if *cur >= len(token){

            return a 
        }
        t = token[*cur].GetType()


        if t == Plus {
            (*cur)++
            b = ParseTerm(token, cur)
            a = NewOpr(a, b, Plus)
        } else {

            return a
        }
    }



    return nil 
}



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


func PrintTree(node TreeNode) {
    fmt.Println("---Print Tree---")
    fmt.Println(node)
    fmt.Println()
}


func SimplifyEquation(tokens []Token) []Token{
    var (
        newTokens []Token
    )
    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == Minus{
            tokens[i].Kind = Plus
            tokens[i + 1].ChangeSign()
        }
    }

    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == Variable || tokens[i].GetType() == VariableNeg {
           if i > 0 && i + 1 < len(tokens) && tokens[i - 1].GetType() == FloatNumber && tokens[i + 1].GetType() == FloatNumber {
                newTokens = append(newTokens, Token{Kind: Multiply})
                newTokens = append(newTokens, tokens[i])
                newTokens = append(newTokens, Token{Kind: Multiply})
           }else if i + 1 < len(tokens) && tokens[i + 1].GetType() == FloatNumber {
                newTokens = append(newTokens, tokens[i])
                newTokens = append(newTokens, Token{Kind: Multiply})
            } else if i > 0 && tokens[i - 1].GetType() == FloatNumber{
                newTokens = append(newTokens, Token{Kind: Multiply})
                newTokens = append(newTokens, tokens[i])
            } else {
                newTokens = append(newTokens, tokens[i])
            }
        }else {
            newTokens = append(newTokens, tokens[i])
        }
    }
    return newTokens
}


func CalculateTree(root TreeNode) Result{
    if root == nil {
        return Result{}
    }

    if node, ok := root.(*OpNode); ok {
        CalculateTree(node.L)
        CalculateTree(node.R)
        node.Res = node.Eval()
    }

    return Result{}
}

// calculate the reduce form

func ReduceForm(c float64){
 	// StackOfVar 
	var (
		mp map[int]float64
	)
	mp = make(map[int]float64)
	for i := 0; i < len(StackOfVar); i++{
		mp[StackOfVar[i].Degree] += StackOfVar[i].Factor 
	}

	fmt.Println("-- Reduced Form --")
	for k,  _ := range(mp){
		keys = append(k, keys)
	}
	sort.
	for k, val := range(mp){
		fmt.Printf("%.2f * X^ %d +", val, k)
	}
	fmt.Println(c)

}

func main() {

	var (
		tokens       []Token
		node        TreeNode
        cur             int
		c				float64
	)
    cur = 0

	if len(os.Args) != 2 {
		Fatal("we need 2 arguments")
	}

	tokens, _ = Parse(os.Args[1])

    tokens = SimplifyEquation(tokens)

	dbgV2(tokens)
    node = ParseExpression(tokens,&cur)
    PrintTree(node)
    CalculateTree(node)

    if res,ok := node.(*OpNode); ok {
        fmt.Println("c = ",res.Res.NumberResult)
		c = res.Res.NumberResult
	}
	ReduceForm(c)
    //    fmt.Println(StackOfVar)
	//dbgV2(tokens)
}
