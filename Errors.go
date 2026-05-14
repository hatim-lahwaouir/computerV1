package main


import (
    "fmt"
    "os"
    "strings"
)


var side bool = true

func ReportErrorSetSide(){
    side = ! side
}


func ReportError(index int , tokens []Token, msg string) {
    var (
        padding int
    )
    padding, _ =  fmt.Printf("Equation : ")

    for i := 0; i < len(tokens); i++ {
        val, _ := fmt.Printf(" %s ", tokens[i])
        if i < index {
            padding += val
        }
    }
    fmt.Println()
    
    pad := strings.Repeat(".", padding + 1) 
    fmt.Printf("%s^\n", pad)
    fmt.Printf("%s\n", msg)

    if side == true {
        fmt.Println("In first side of equation")
    } else {
        fmt.Println("In second side of equation")
    }
    os.Exit(1)
}

func ReportErrortringString(index int , tokens string, msg string) {
    var (
        padding int
    )
    padding, _ =  fmt.Printf("Equation : ")

    for i := 0; i < len(tokens); i++ {
        val, _ := fmt.Printf("%c", tokens[i])
        if i < index {
            padding += val
        }
    }
    fmt.Println()
    
    pad := strings.Repeat(".", padding) 
    fmt.Printf("%s^\n", pad)
    fmt.Printf("%s\n", msg)

    if side == true {
        fmt.Println("In first side of equation")
    } else {
        fmt.Println("In second side of equation")
    }
    os.Exit(1)
}
