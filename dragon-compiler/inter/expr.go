package inter

import (
	"lexer"
)

type Expr struct {
	Node      *Node
	token     *lexer.Token
	expr_type *Type
}

func NewExpr(line uint32, token *lexer.Token, expr_type *Type) *Expr {
	expr := &Expr{
		Node:      NewNode(line),
		token:     token,
		expr_type: expr_type,
	}

	return expr
}

func (e *Expr) Errors(s string) error {
	return e.Node.Errors(s)
}

func (e *Expr) NewLabel() uint32 {
	return e.Node.NewLabel()
}

func (e *Expr) EmitLabel(i uint32) {
	e.Node.EmitLabel(i)
}

func (e *Expr) Emit(code string) {
	e.Node.Emit(code)
}

func (e *Expr) Gen() ExprInterface {
	return e
}

func (e *Expr) Reduce() ExprInterface {
	return e
}

func (e *Expr) ToString() string {
	return e.token.ToString()
}

func (e *Expr) Type() *Type {
	return e.expr_type
}
