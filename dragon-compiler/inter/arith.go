package inter

import (
	"lexer"
)


type Arith struct {
	op        *Op
	line      uint32
	token     *lexer.Token
	expr1     ExprInterface
	expr2     ExprInterface
	expr_type *Type
}

func NewArith(line uint32, token *lexer.Token, expr1 ExprInterface,
	expr2 ExprInterface) (*Arith, error) {
	/*
		根据expr1 , expr2对应的类型进行提升
	*/
	expr_type := MaxType(expr1.Type(), expr2.Type())
	arith := &Arith{
		op:        NewOp(line, token, expr_type),
		line:      line,
		token:     token,
		expr1:     expr1,
		expr2:     expr2,
		expr_type: expr_type,
	}

	arith.op.child = arith

	if expr_type == nil {
		//表达式的类型不能为空
		return nil, arith.op.Errors("type error")
	}

	return arith, nil
}

func (a *Arith) Errors(s string) error {
	return a.op.Errors(s)
}

func (a *Arith) NewLabel() uint32 {
	return a.op.NewLabel()
}

func (a *Arith) EmitLabel(i uint32) {
	a.op.EmitLabel(i)
}

func (a *Arith) Emit(code string) {
	a.op.Emit(code)
}

func (a *Arith) Gen() ExprInterface {
	/*
		我们可能会遇到复杂的组合表达式例如 (a+b) + (c+d)，
		于是expr1 对应a+b, expr2 对应c+d，
		此时节点生成中间代码时，需要先让expr1和expr2生成代码
	*/
	arith, _ := NewArith(a.line, a.token, a.expr1.Reduce(), a.expr2.Reduce())
	return arith
}

func (a *Arith) Reduce() ExprInterface {
	return a.op.Reduce()
}

func (a *Arith) ToString() string {
	return a.expr1.ToString() + " " + a.token.ToString() + " " + a.expr2.ToString()
}

func (a *Arith) Type() *Type {
	return a.expr_type
}
