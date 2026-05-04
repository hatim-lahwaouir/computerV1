package main

import (
    "os"
    "fmt"
    "strconv"
    "unicode"

)

type Tokens int

const (
    Plus = iota
    Minus
    Multiply
    Caret
    Intnumber  
    FloatNumber 
    Variable 
)



var TokenState = map[string]Tokens{
    "+": Plus,
    "-": Minus,
    "*": Multiply,
    "^": Caret,
    "X": Variable,
    "x": Variable,
    
}

var TokenStr = map[Tokens]string{
    Plus : "+",
    Minus: "-",
    Multiply: "*",
    Caret : "^",
    Variable: "X",
}



// error function

func Fatal(msg string) {


    fmt.Printf("error: %s\n", msg)

    os.Exit(1)
}


func main(){

    var (
        input string
        tokens []Tokens
        numbers []int
    )

    if len(os.Args) != 2{
        Fatal("we need 2 arguments")
    }
    
    input = os.Args[1]

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
            )
            tokens = append(tokens, Intnumber)
            s = i
            for ; i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.'); i++{

            }
            e = i 
            val, err := strconv.Atoi(input[s:e])

            if err != nil {
                fmt.Println(input[s:e], s, e)
                Fatal("error parssing number")
            }
            numbers = append(numbers, val)
            i--;
        }
    }

    // dbg

    j := 0
    for i := 0; i < len(tokens); i++{
        if val, ok := TokenStr[tokens[i]]; ok  == true{
            fmt.Printf("%s ", val)
        }else {
            fmt.Printf("%d", numbers[j])
            j++
        }
    }
}
