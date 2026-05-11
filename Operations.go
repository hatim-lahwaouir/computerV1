package main


import (
    "math"
    "fmt"
)


var StackOfVar []Var 


type Result struct {
    IsVar bool
    VarResult Var
    NumberResult float64 
}


func (l Result) Add(r Result) Result {
    if l.IsVar == false && r.IsVar == false{

        return Result{IsVar: false, NumberResult: l.NumberResult + r.NumberResult}
    } else {

        if l.IsVar && r.IsVar && r.VarResult.Degree ==  l.VarResult.Degree{
            fmt.Println(">>", r, l)
            return l.VarResult.Add(r.VarResult) 
        }else if l.IsVar && r.IsVar && r.VarResult.Degree !=  l.VarResult.Degree{

            StackOfVar = append(StackOfVar,r.VarResult)
            StackOfVar = append(StackOfVar,l.VarResult)
        }else if l.IsVar {

            StackOfVar = append(StackOfVar,l.VarResult)
            return r
        }else {


            StackOfVar = append(StackOfVar,r.VarResult)
            return l
        }

            fmt.Println("here", l, r)

    }
    return Result{}
}

func (l Result) Multiplication(r Result) Result {
    if l.IsVar == false && r.IsVar == false{
        return Result{IsVar: false, NumberResult: l.NumberResult * r.NumberResult}
    } else {
        if l.IsVar && r.IsVar {
            return l.VarResult.Multiplication(r.VarResult) 
        }else if l.IsVar {
            return l.VarResult.MultiplyByFactor(r.NumberResult)
        }else {
            return r.VarResult.MultiplyByFactor(l.NumberResult)
        }
    }
    return Result{}
}


func (l Result) CaretOperation(r Result) Result {
    if l.IsVar == false && r.IsVar == false{
        return Result{IsVar: false, NumberResult: math.Pow(l.NumberResult , r.NumberResult)}
    } else {
        if l.IsVar {
            fmt.Println(">>>> ^ ", l, r)
            return l.VarResult.SetDegree(int(r.NumberResult)) 
        }
    }
    return Result{}
}


type Var struct {
    Degree int 
    Factor float64 
}

func NewVar(dg int, Factor float64) TreeNode{
    return &Var{ Degree: dg,Factor: Factor}
}

func (o *Var) Eval() Result{
    return Result{IsVar: true, VarResult: *o}
}

func (l Var) Add(r Var) Result{
    return Result {IsVar: true, VarResult: Var{r.Degree, l.Factor + r.Factor}}
}

func (l Var) Multiplication(r Var) Result{
    return Result {IsVar: true, VarResult: Var{r.Degree + l.Degree, l.Factor * r.Factor}}
}

func (l Var) SetDegree(degree int) Result{
    return Result {IsVar: true, VarResult: Var{degree, l.Factor}}
}

func (l Var) MultiplyByFactor(factor float64) Result{
    return Result {IsVar: true, VarResult: Var{l.Degree, l.Factor * factor}}
}



func (o Var) String() string{
    if o.Factor <  0 {
        return fmt.Sprintf("(%.2f * x^%d)",o.Factor ,o.Degree)
    }
    return fmt.Sprintf("(%.2f * x^%d)", o.Factor, o.Degree)
}

func (o Var) GetKind() TokenType{
    return Variable 
}


type Value struct {
    Val float64 
}

func (o *Value) Eval() Result{
    return Result{IsVar: false, NumberResult : o.Val}
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

func NewVal(val float64) TreeNode{
    return &Value{Val: val} 
}



type TreeNode interface {
    Eval() Result 
    GetVal() Result 
    GetKind() TokenType
}




type OpNode struct {
    Op    TokenType 
    Res   Result 
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

func (o *OpNode) Eval() Result{
    switch o.Op {
    case Plus:
        return o.L.GetVal().Add(o.R.GetVal())
    case Multiply:
        return o.L.GetVal().Multiplication(o.R.GetVal())
    case Caret:  
        return o.L.GetVal().CaretOperation(o.R.Eval())
    }

    return Result{} 
}

func (o OpNode) GetVal() Result{
    return o.Res
}

func (o Value) GetVal() Result{
    return Result{IsVar: false, NumberResult : o.Val}
}

func (o Var) GetVal() Result{
    return Result{IsVar: true, VarResult: o}
}



func NewOpr(a TreeNode, b TreeNode, t TokenType ) TreeNode{
    return &OpNode{Op: t, L: a, R: b} 
}






func (o OpNode) GetKind() TokenType{
    return o.Op
}




