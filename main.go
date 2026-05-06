package main

import (
    "os"
    "fmt"
    "strconv"
    "unicode"
    "errors"

)

type Tokens int

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
)



var TokenState = map[string]Tokens{
    "+": Plus,
    "-": Minus,
    "*": Multiply,
    "^": Caret,
    "X": Variable,
    "x": Variable,
    "x^2": VariableDegree2,
    "x^1": VariableDegree1,
    
}

var TokenStr = map[Tokens]string{
    Plus : "+",
    Minus: "-",
    Multiply: "*",
    Caret : "^",
    Variable: "X",
    VariableDegree2: "x^2",
    VariableDegree1: "x^1",
}



// error function

func Fatal(msg string) {


    fmt.Printf("error: %s\n", msg)

    os.Exit(1)
}
// helpers 

func isOperation(t Tokens) bool {
    if t != Plus && t != Minus && t != Multiply && Caret != t {
        return false
    } 

    return true
}



func SimplifyVariables1 (tokens []Tokens, numbers []int, floatNumbers[]float32) ([]Tokens, []int, []float32, error) {
    var (
        newTokens []Tokens
        newNumbers []int
        numberIndex int
    )


    numberIndex = 0
    
   


    for i := 0; i < len(tokens); i++ {

        if tokens[i] == Variable {
            if i + 1 == len(tokens) {
                newTokens = append(newTokens, VariableDegree1)
            } else {

                if isOperation(tokens[i + 1]) == false {
                    return nil,nil,nil ,errors.New("Invalid equation operations must be after the variable")
                }

                if tokens[i + 1] == Caret {
                    if i + 2 >= len(tokens) || tokens[i + 2] != IntNumber {
                        return nil, nil, nil,errors.New("Invalid equation no number was provided after this \"^\"")
                    }

                    switch numbers[numberIndex] {
                        case 0:
                            newNumbers = append(newNumbers, 1)
                            newTokens = append(newTokens, IntNumber)
                        case 1:
                            newTokens = append(newTokens, VariableDegree1)
                        case 2:
                            newTokens = append(newTokens, VariableDegree2)
                        default:
                            return nil,nil,nil, errors.New("we only solve second degree equations")
                    }
                    numberIndex++
                    i+=2 
                } else {
                            newTokens = append(newTokens, VariableDegree1)
                }
            }

        } else {
            if tokens[i] == IntNumber {
                newNumbers = append(newNumbers, numbers[numberIndex])
                numberIndex++
            }
            newTokens = append(newTokens, tokens[i])
        }
    }


    return newTokens, newNumbers, floatNumbers, nil
}


func SimplifyVariables2 (tokens []Tokens) ([]Tokens, error) {

    var (
            newTokens []Tokens
        )

        for i := 0; i < len(tokens); i++ {
                if tokens[i] == Empty{
                    continue
                }
        
                
                if tokens[i] == VariableDegree1 || tokens[i] ==  VariableDegree2 {
                    var (
                        oper bool
                        degree int

                    )
                    oper = false 
                    degree = 0 



                    if i + 1 >= len(tokens) {
                        newTokens = append(newTokens, tokens[i])
                        break 
                    }
                    j := i
                    for ; j < len(tokens); j++{
                        if oper {
                            if tokens[j] != Multiply {
                                break
                            }
                        } else {
                            if degree == 1  && (tokens[j] == VariableDegree2 ||  tokens[j] == VariableDegree1){
                                tokens[j- 1] = Empty 
                            }
                            
                            if tokens[j] == VariableDegree1 {
                                tokens[j] = Empty 
                                degree++
                            } else if  tokens[j] == VariableDegree2 {
                                tokens[j] = Empty 
                                degree += 2
                            }

                            
                        }


                        
                        oper = !oper
                    }
                    //tokens[j - 2] = Empty

                    if degree == 1 {
                            newTokens = append(newTokens, VariableDegree1)
                    } else if degree == 2 {
                            newTokens = append(newTokens, VariableDegree2)
                    }else {
                            return nil, errors.New("we only solve second degree equations")
                    }
                    
                }else {
                    newTokens = append(newTokens, tokens[i])
                }
        }
        return newTokens, nil

}

func Parse(input string) ([]Tokens, []int,[]float32, error) {
    var (
        tokens []Tokens
        numbers []int
        floatNumbers []float32
    )


    for i := 0; i < len(input); i++ {

        if unicode.IsSpace(rune(input[i])){
            continue
        }
        
        if val , ok := TokenState[string(input[i])]; ok == true {
            tokens = append(tokens, val)

        }else {
            // treating the number
            var (
                s int
                e int
                isFloat bool
            )

            isFloat = false
            s = i
            for ; i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.'); i++{
                if input[i] == '.'{
                    isFloat = true

                }

            }
            e = i 
            if isFloat == true {
                tokens = append(tokens, FloatNumber)
                val, err := strconv.ParseFloat(input[s:e], 32)

                if err != nil {
                    fmt.Println(input[s:e], s, e)
                    Fatal("error parssing number")
                }
                floatNumbers = append(floatNumbers, float32(val))

            }else {

                tokens = append(tokens, IntNumber)
                val, err := strconv.Atoi(input[s:e])
                if err != nil {
                    fmt.Println(input[s:e], s, e)
                    Fatal("error parssing number")
                }
                numbers = append(numbers, val)
            }
            i--;
        }
    }

    return  tokens, numbers, floatNumbers ,nil
}


func dbg (tokens []Tokens, numbers []int, floatNumbers[]float32){

    j := 0
    k := 0 
    for i := 0; i < len(tokens); i++{
        val, ok := TokenStr[tokens[i]];
        if  ok  == true{
            fmt.Printf(" %s ", val)
        }else if tokens[i] ==  IntNumber{
            fmt.Printf(" %d ", numbers[j])
            j++
        } else {
            fmt.Printf(" %.2f ", floatNumbers[k])
            k++
        }
    }

}


func main(){

    var (
        tokens []Tokens
        numbers []int
        floatNumbers []float32
    )

    if len(os.Args) != 2{
        Fatal("we need 2 arguments")
    }
    
    tokens, numbers,floatNumbers, _ = Parse(os.Args[1])

    tokens, numbers, floatNumbers, err := SimplifyVariables1(tokens, numbers, floatNumbers)

    if err != nil {
        Fatal(err.Error())
    }

    tokens, err = SimplifyVariables2(tokens)
    if err != nil {
        Fatal(err.Error())
    }
    dbg(tokens, numbers, floatNumbers)
    

}
