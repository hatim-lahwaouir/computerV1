package main

import (
	"fmt"
	"os"
    "strings"
	"strconv"
	"unicode"
	"sort"
)






// F = X | Number | X ^ CARET  
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
				val, err := strconv.ParseFloat(input[s:e], 64)

				if err != nil {
                    ReportError(len(tokens) - 1, tokens, "Invalid number " +  input[s:e])
				}
				tokens = append(tokens, Token{Kind: FloatNumber, Value: float64(val) })
			} else {
				
				val, err := strconv.ParseFloat(input[s:e], 64)
				if err != nil {
                    ReportError(len(tokens) - 1, tokens, "Invalid number " +  input[s:e] )
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



func SimplifyEquation(tokens []Token) []Token{
    var (
        newTokens []Token
    )

    // remove + from the first token
    if tokens[0].Kind == Plus{
        tokens[0] = Token{Kind: Empty}
    }
    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == Minus{
            tokens[i].Kind = Plus
            if  i + 1 < len(tokens) {
                tokens[i + 1].ChangeSign()
            }
            if i == 0 && tokens[i].Kind == Plus{
                tokens[0] = Token{Kind: Empty}
            }
        }
    }

    for i := 0; i < len(tokens); i++{
        if tokens[i].GetType() == Empty {
            continue
        }
        if tokens[i].GetType() == Variable || tokens[i].GetType() == VariableNeg {
            if i > 0 && tokens[i - 1].GetType() == FloatNumber{
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


func CalculateTree(root TreeNode) {
    if root == nil {
        return 
    }

    if node, ok := root.(*OpNode); ok {
        CalculateTree(node.L)
        CalculateTree(node.R)
        node.Res = node.Eval()
        return
    }
    root.Eval()
}

// calculate the reduce form

func ReduceForm(side1 []Var, side2 []Var,  c float64) (map[int]float64, []int){
    var (
		mp map[int]float64
        keys   []int
	)

	mp = make(map[int]float64)
	for i := 0; i < len(side1); i++{
		mp[side1[i].Degree] += side1[i].Factor 
	}

    for i := 0; i < len(side2); i++{
		mp[side2[i].Degree] -= side2[i].Factor 
	}

	fmt.Println("-- Reduced Form --")

	for k,  _ := range(mp){
		keys = append(keys,k)
	}
    sort.Ints(keys)



    fmt.Printf("(%.2f * X ^ 0 ) + ", c )
	for i, k := range(keys){
        if mp[k] != 0 {
	        fmt.Printf("(%.2f * X^ %d)", mp[k], k)
        }
        if i + 1 != len(keys){
            fmt.Printf(" + ")
        }
	}
    fmt.Println(" = 0")
    return mp, keys
}

// Validation 

func Validation(tokens []Token) int {
    var (
        IsOpr  bool
    )


    IsOpr = false

    if len(tokens) == 0 {
        return 0  
    }

    for i := 0 ; i < len(tokens) ; i++ {
        if IsOpr {
            if IsOperation(tokens[i]) == false{
                return i
            }
        } else {
            if IsVariable(tokens[i]) == false && tokens[i].GetType() != FloatNumber {
                return i
            }
        }
        IsOpr = ! IsOpr
    }

    // last token must be a number
    if IsOpr  == false {
                return  len(tokens) - 1 
    }

    for i := 0; i < len(tokens); i++ {
        if tokens[i].GetType() == Caret{
            if tokens[i + 1].Value !=  float64(int64(tokens[i + 1].Value)) {
                return i + 1
            }
            if IsVariable(tokens[i + 1]) {
                return i + 1
            }
        }
    }


    return -1
}


func ParseSide(input string) float64 {
    var (
		tokens       []Token
		node        TreeNode
        c           Result
        cur             int
	)
    cur = 0


	tokens, _ = Parse(input)
    tokens = SimplifyEquation(tokens)
    if index := Validation(tokens); index != -1 {
        ReportError(index, tokens, "Invalid Operation (program change - to + | pay attention to this )")
    }
    node = ParseExpression(tokens,&cur)
    CalculateTree(node)

    c = node.GetVal() 
    if c.IsVar {
            StackOfVar = append(StackOfVar, c.VarResult)
            return 0
    }
    return c.NumberResult 
}


func PolynomialDegree(keys []int, mp map[int]float64) int {
	fmt.Println("-- PolynomialDegree --")

    for i := len(keys) - 1; i >= 0; i--{

        if mp[keys[i]] != 0 {
            fmt.Println(keys[i])
            return keys[i]
        }
    }
    fmt.Println(0)
    return 0
}
func findSolution(equationDegree int ,  mp map[int]float64,  c float64){
    var (
        delta float64
        a float64
        b float64

    )
    

    if equationDegree == 1 {
        if c == 0 {
            fmt.Printf("the solution is 0\n")
        } else {
            fmt.Printf("the solution is %.3f\n", -c / mp[1] ) 
        }
    } else if  equationDegree == 2 {
        a = mp[2]
        b =  mp[1]
        delta =  (b * b) - (4 * a * c)
        if delta > 0 {
            fmt.Println("Discriminant is strictly positive, the two solutions are:")
            fmt.Printf("%.3f\n",(-b + Sqrt(delta)) / (2 * a))
            fmt.Printf("%.3f\n",(-b - Sqrt(delta)) / (2 * a))

        } else if delta == 0 {
            fmt.Println("Discriminant is equalt to zero, the only solutions is:")
            fmt.Println(-b / (2 * a))

        } else if delta < 0 {
            fmt.Println("Discriminant is strictly negative, the two complex solutions are:")
            fmt.Printf("%.3f + %.3fi\n", -b/ (2 * a), Sqrt(- delta) / (2 * a) )
            fmt.Printf("%.3f - %.3fi\n",-b/ (2 * a),Sqrt(-delta) / (2 * a) )
        }
    }
}

func StartParsing(input string){
    var (
        Sides []string
        c1 float64
        c2 float64
        VarOfSide1 []Var   
        VarOfSide2 []Var   
		mp map[int]float64
        keys   []int
        equationDegree int
    )

    Sides = strings.Split(input, "=")

    if len(Sides) != 2 {
        Fatal("invalide equation '=', either more then one was provided or non was provided ")
    }
    c1 = ParseSide(Sides[0])
    // reset the stack for Xs in the other side
    VarOfSide1 = StackOfVar
    StackOfVar = []Var{} 
    // change side of reporting errors
    ReportErrorSetSide()
    c2 = ParseSide(Sides[1])
    VarOfSide2 = StackOfVar

    
    mp, keys = ReduceForm(VarOfSide1, VarOfSide2, c1 - c2)

    equationDegree = PolynomialDegree(keys, mp)

    if equationDegree > 2 {
        fmt.Println("The polynomial degree is strictly greater than 2, I can't solve.")
        return

    }
    if equationDegree == 0 {
        if c1 - c2 == 0 {
            fmt.Println("Any real number is a solution.")
        } else {
            fmt.Println("No solution.")
        }
    } else {
        findSolution(equationDegree,  mp,  c1 - c2)

    }
    
}



func main() {

	if len(os.Args) != 2 {
		Fatal("we need 2 arguments")
	}
    StartParsing(os.Args[1])

	//    fmt.Println(StackOfVar)
	//dbgV2(tokens)
}
