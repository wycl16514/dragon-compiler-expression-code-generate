package inter

type Seq struct {
	Node  *Node
	stmt1 StmtInterface
	stmt2 StmtInterface
}

func NewSeq(line uint32, stmt1 StmtInterface, stmt2 StmtInterface) *Seq {
	return &Seq{
		Node:  NewNode(line),
		stmt1: stmt1,
		stmt2: stmt2,
	}
}

func (s *Seq) Errors(str string) error {
	return s.Node.Errors(str)
}

func (s *Seq) NewLabel() uint32 {
	return s.Node.NewLabel()
}

func (s *Seq) EmitLabel(i uint32) {
	s.Node.EmitLabel(i)
}

func (s *Seq) Emit(code string) {
	s.Node.Emit(code)
}

func (s *Seq) Gen(start uint32, end uint32) {
	_, ok1 := s.stmt1.(*Stmt)
	_, ok2 := s.stmt2.(*Stmt)
	if ok1 {
		s.stmt2.Gen(start, end)
	} else if ok2 {
		s.stmt1.Gen(start, end)
	} else {
		label := s.Node.NewLabel()
		s.stmt1.Gen(start, label)
		s.Node.EmitLabel(label)
		s.stmt2.Gen(label, end)
	}
}
