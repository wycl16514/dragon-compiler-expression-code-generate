package inter

import (
	"lexer"
)

type Op struct {
	expr      *Expr
	child     ExprInterface
	line      uint32
	expr_type *Type
}

func NewOp(line uint32, token *lexer.Token, expr_type *Type) *Op {
	op := &Op{
		expr:      NewExpr(line, token, expr_type),
		child:     nil,
		line:      line,
		expr_type: expr_type,
	}

	return op
}

func (o *Op) Errors(s string) error {
	return o.expr.Errors(s)
}

func (o *Op) NewLabel() uint32 {
	return o.expr.NewLabel()
}

func (o *Op) EmitLabel(i uint32) {
	o.expr.EmitLabel(i)
}

func (o *Op) Emit(code string) {
	o.expr.Emit(code)
}

func (o *Op) Gen() ExprInterface {
	return o
}

func (o *Op) Reduce() ExprInterface {
	if o.child != nil {
		/*调用子节点的Gen函数，让子节点先生成中间代码,
			子节点生成中间代码后会返回一个Expr节点，然后这里将返回的节点赋值给
			一个临时寄存器变量

			具体逻辑为当编译器遇到语句 a + b 就会生成Op节点,
		    那么a + b对应一个Arith节点，它对应child对象，
			执行child.Gen()会生成中间代码对应的字符串:
			a + b
		    接下来我们创建一个临时寄存器变量例如t2,然后生成语句
			t2 = a + b
		*/
		x := o.child.Gen()
		t := NewTemp(o.line, o.expr_type)
		o.expr.Emit(t.ToString() + " = " + x.ToString())
		return t
	}

	return nil
}

func (o *Op) ToString() string {
	return o.expr.ToString()
}

func (o *Op) Type() *Type {
	return o.expr.Type()
}
