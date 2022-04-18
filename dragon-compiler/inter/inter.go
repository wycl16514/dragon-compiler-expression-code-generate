package inter

type NodeInterface interface {
	Errors(s string) error
	NewLabel() uint32
	EmitLabel(i uint32)
	Emit(code string)
}

type ExprInterface interface {
	NodeInterface
	Gen() ExprInterface
	Reduce() ExprInterface
	Type() *Type
	ToString() string
}

type StmtInterface interface {
	NodeInterface
	//start, end对应语句的起始标签和结束标签号码
	Gen(start uint32, end uint32)
}
