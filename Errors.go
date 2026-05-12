package main


import (
    "fmt"
    "os"
    "strings"
)



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

    os.Exit(1)
}
