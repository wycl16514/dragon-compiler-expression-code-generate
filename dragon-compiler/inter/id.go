package inter

import (
	"lexer"
)

type ID struct {
	/*
		该节点没有实现Gen,Reduce()，这意味着编译器遇到语句:"a;","b;"等时会直接越过
		不生成任何中间代码
	*/
	expr   *Expr
	Offset uint32 //相对偏移地址，用于生成中间代码，
}

func NewID(line uint32, token *lexer.Token, expr_type *Type) *ID {
	id := &ID{
		expr: NewExpr(line, token, expr_type),
	}

	return id
}

func (i *ID) Errors(s string) error {
	return i.expr.Errors(s)
}

func (i *ID) NewLabel() uint32 {
	return i.expr.NewLabel()
}

func (i *ID) EmitLabel(l uint32) {
	i.expr.EmitLabel(l)
}

func (i *ID) Emit(code string) {
	i.expr.Emit(code)
}

func (i *ID) Gen() ExprInterface {
	return i
}

func (i *ID) Reduce() ExprInterface {
	return i
}

func (i *ID) Type() *Type {
	return i.expr.Type()
}

func (i *ID) ToString() string {
	return i.expr.ToString()
}
