package inter

/*
Set 节点对应 c = a+b,因此它包含两部分，分别是左边的ID节点和右边的Op节点
*/
type Set struct {
	id   ExprInterface
	expr ExprInterface
}

func checkType(p1 *Type, p2 *Type) *Type {
	//c = a + b , c的类型会转换为右边a+b的类型
	if Numberic(p1) && Numberic(p2) {
		return p2
	} else if p1.Lexeme == "bool" && p2.Lexeme == "bool" {
		return p2
	}

	return nil
}

func NewSet(id ExprInterface, expr ExprInterface) (*Set, error) {
	if checkType(id.Type(), expr.Type()) == nil {
		return nil, id.Errors("type error")
	} else {
		return &Set{
			id:   id,
			expr: expr,
		}, nil
	}
}

func (s *Set) Errors(str string) error {
	return s.id.Errors(str)
}

func (s *Set) NewLabel() uint32 {
	return s.id.NewLabel()
}

func (s *Set) EmitLabel(i uint32) {
	s.id.EmitLabel(i)
}

func (s *Set) Emit(code string) {
	s.id.Emit(code)
}

func (s *Set) Gen() ExprInterface {
	s.expr = s.expr.Gen()
	s.Emit(s.id.ToString() + " = " + s.expr.ToString())
	return s.id
}

func (s *Set) Reduce() ExprInterface {
	return s.id.Reduce()
}

func (s *Set) Type() *Type {
	return s.id.Type()
}

func (s *Set) ToString() string {
	return s.id.ToString()
}
