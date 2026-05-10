package main


import (
    "math"
    "fmt"

)


type Var struct {
    Degree int 
    Negatif bool 
}

func NewVar(negatif bool, dg int) TreeNode{
    return Var{ Degree: dg, Negatif: negatif}

}

func (o Var) Eval() float64{
    return 0 
}

func (o Var) String() string{
    if o.Negatif == true{
        return fmt.Sprintf("(-x)^%d", o.Degree)
    }

    return fmt.Sprintf("(x)^%d", o.Degree)
}

func (o Var) GetKind() TokenType{
    return Variable 
}


type Value struct {
    Val float64 
}

func (o Value) Eval() float64{
    return o.Val
}

func (o Value) String() string{
    if o.Val < 0 {
        return fmt.Sprintf("(%.2f)", o.Val ) 
    }

    return fmt.Sprintf("%.2f", o.Val ) 
}

func (o Value) GetKind() TokenType{
    return FloatNumber 
}

type TreeNode interface {
    Eval() float64 
    GetKind() TokenType
}




type OpNode struct {
    Op    TokenType 
    L TreeNode 
    R TreeNode 
}

func (o OpNode) String() string {
    if o.L != nil &&  o.R != nil{
        return fmt.Sprintf("%s %s %s", o.L,TokenStr[o.Op], o.R)
    } else if o.L != nil {
        return fmt.Sprintf("%s %s", o.L,TokenStr[o.Op])
    } else {
        return fmt.Sprintf("%s %s",TokenStr[o.Op], o.R)
    }
}

func (o OpNode) Eval() float64{
    switch o.Op {
    case Plus:
        return o.L.Eval() + o.R.Eval()
    case Multiply:
        return o.L.Eval() * o.R.Eval()
    case Caret:  
        return math.Pow(o.L.Eval() , o.R.Eval())
    case Minus:  
        return o.L.Eval() - o.R.Eval()
    default:
        return 0
    }
}

func NewOpr(a TreeNode, b TreeNode, t TokenType ) TreeNode{
    return OpNode{Op: t, L: a, R: b} 
}





func NewVal(val float64) TreeNode{
    return Value{Val: val} 
}


func (o OpNode) GetKind() TokenType{
    return o.Op
}




