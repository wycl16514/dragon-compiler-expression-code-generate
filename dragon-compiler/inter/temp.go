package inter

import (
	"lexer"
	"strconv"
)

/*
Temp节点表示中间代码中的临时寄存器变量
*/
type Temp struct {
	expr   *Expr
	number uint32
}

var count uint32

func NewTemp(line uint32, expr_type *Type) *Temp {
	token := lexer.NewToken(lexer.TEMP)
	count = count + 1

	temp := &Temp{
		expr:   NewExpr(line, &token, expr_type),
		number: count,
	}

	return temp
}

func (t *Temp) Errors(s string) error {
	return t.expr.Errors(s)
}

func (t *Temp) NewLabel() uint32 {
	return t.expr.NewLabel()
}

func (t *Temp) EmitLabel(i uint32) {
	t.expr.EmitLabel(i)
}

func (t *Temp) Emit(code string) {
	t.expr.Emit(code)
}

func (t *Temp) Gen() ExprInterface {
	return t
}

func (t *Temp) Reduce() ExprInterface {
	return t
}

func (t *Temp) ToString() string {
	return "t" + strconv.FormatUint(uint64(t.number), 10)
}

func (t *Temp) Type() *Type {
	return t.expr.Type()
}
