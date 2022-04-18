package inter

type Stmt struct {
	Node      *Node
	after     uint32        //用于生成跳转标签
	enclosing StmtInterface //用于break语句
}

func NewStmt(line uint32) *Stmt {
	return &Stmt{
		Node:      NewNode(line),
		after:     0,
		enclosing: nil,
	}
}

func (s *Stmt) Errors(str string) error {
	return s.Node.Errors(str)
}

func (s *Stmt) NewLabel() uint32 {
	return s.Node.NewLabel()
}

func (s *Stmt) EmitLabel(i uint32) {
	s.Node.EmitLabel(i)
}

func (s *Stmt) Emit(code string) {
	s.Node.Emit(code)
}

func (s *Stmt) Gen(start uint32, end uint32) {
	//这里需要子节点来实现

}
