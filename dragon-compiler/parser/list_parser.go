package simple_parser

import (
	"errors"
	"fmt"
	"inter"
	"lexer"
)

type SimpleParser struct {
	lexer        lexer.Lexer
	top          *Env
	saved        *Env
	cur_tok      lexer.Token
	used_storage uint32 //当前用于存储变量的内存字节数
}

func NewSimpleParser(lexer lexer.Lexer) *SimpleParser {
	return &SimpleParser{
		lexer: lexer,
		top:   nil,
		saved: nil,
	}
}

func (s *SimpleParser) Parse() {
	s.program()
}

func (s *SimpleParser) program() {
	s.top = nil
	//stmt 其实是seq所形成的队列的头结点
	stmt := s.block()

	begin := stmt.NewLabel()
	after := stmt.NewLabel()
	stmt.EmitLabel(begin)
	stmt.Gen(begin, after)
	stmt.EmitLabel(after)

}

func (s *SimpleParser) matchLexeme(str string) error {
	if s.lexer.Lexeme == str {
		return nil
	}

	err_s := fmt.Sprintf("error token , expected:%s , got:%s", str, s.lexer.Lexeme)
	return errors.New(err_s)
}

func (s *SimpleParser) matchTag(tag lexer.Tag) error {
	if s.cur_tok.Tag == tag {
		return nil
	}

	err_s := fmt.Sprintf("error tag, expected:%d, got %d", tag, s.cur_tok.Tag)
	return errors.New(err_s)
}

func (s *SimpleParser) move_backward() {
	s.lexer.ReverseScan()
}

func (s *SimpleParser) move_forward() error {
	var err error
	s.cur_tok, err = s.lexer.Scan()
	return err
}

func (s *SimpleParser) block() inter.StmtInterface {
	// block -> "{" decls statms "}"
	err := s.move_forward()
	if err != nil {
		panic(err)
	}

	err = s.matchLexeme("{")
	if err != nil {
		panic(err)
	}

	err = s.move_forward()
	if err != nil {
		panic(err)
	}

	s.saved = s.top
	s.top = NewEnv(s.top)
	err = s.decls()
	if err != nil {
		panic(err)
	}

	stmt := s.stmts()
	if err != nil {
		panic(err)
	}

	err = s.matchLexeme("}")
	if err != nil {
		panic(err)
	}

	s.top = s.saved
	return stmt
}

func (s *SimpleParser) decls() error {
	/*
		decls -> decls decl | ε
		decls 表示由零个或多个decl组成，decl对应语句为:
		int a; float b; char c;等，其中int, float, char对应的标号为BASIC,
		在进入到这里时我们并不知道要解析多少个decl,一个处理办法就是判断当前读到的字符串标号，
		如果当前读到了BASIC标号，那意味着我们遇到了一个decl对应的声明语句，于是就执行decl对应的语法
		解析，完成后我们再次判断接下来读到的是不是还是BASIC标号，如果是的话继续进行decl解析，
		由此我们可以破除左递归
	*/
	for s.cur_tok.Tag == lexer.BASIC {
		err := s.decl()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SimpleParser) getType() (*inter.Type, error) {
	err := s.matchTag(lexer.BASIC)
	if err != nil {
		return nil, err
	}

	width := uint32(4)
	switch s.lexer.Lexeme {
	case "int":
		width = 4
	case "float":
		width = 8
	case "char":
		width = 1
	}

	p := inter.NewType(s.lexer.Lexeme, lexer.BASIC, width)
	s.used_storage = s.used_storage + width
	return p, nil
}

func (s *SimpleParser) decl() error {
	p, err := s.getType()
	if err != nil {
		return err
	}

	err = s.move_forward()
	if err != nil {
		return err
	}
	//这里必须复制，因为s.cur_tok会不断变化因此不能直接传入s.cur_tok
	tok := lexer.NewTokenWithString(s.cur_tok.Tag, s.lexer.Lexeme)
	id := inter.NewID(s.lexer.Line, tok, p)

	sym := NewSymbol(id, p)
	s.top.Put(s.lexer.Lexeme, sym)

	err = s.move_forward()
	if err != nil {
		return err
	}

	err = s.matchLexeme(";")
	if err != nil {
		return err
	}

	err = s.move_forward()
	return err
}

func (s *SimpleParser) stmts() inter.StmtInterface {
	if s.matchLexeme("}") == nil {
		return inter.NewStmt(s.lexer.Line)
	}

	//注意这里，seq节点通过递归形成了一个链表
	return inter.NewSeq(s.lexer.Line, s.stmt(), s.stmts())
}

func (s *SimpleParser) stmt() inter.StmtInterface {
	return s.expression()
}

func (s *SimpleParser) expression() inter.StmtInterface {
	if s.matchTag(lexer.ID) == nil {
		s.move_forward()
		if s.matchTag(lexer.ASSIGN_OPERATOR) == nil {
			s.move_backward()
			s.move_backward() //回退到变量名
			return s.assign()
		}
		s.move_backward()
	}

	expression := inter.NewExpression(s.lexer.Line, s.expr())
	return expression
}

func (s *SimpleParser) assign() inter.StmtInterface {
	s.move_forward()
	sym := s.top.Get(s.lexer.Lexeme)
	if sym == nil {
		err_s := fmt.Sprintf("undefined variable with name: %s", s.lexer.Lexeme)
		err := errors.New(err_s)
		panic(err)
	}

	s.move_forward() //读取=
	s.move_forward() //读取 = 后面的字符串
	expr := s.expr()
	set, err := inter.NewSet(sym.id, expr)
	if err != nil {
		panic(err)
	}
	err = s.matchLexeme(";")
	if err != nil {
		panic(err)
	}
	s.move_forward()
	expression := inter.NewExpression(s.lexer.Line, set)
	return expression
}

func (s *SimpleParser) expr() inter.ExprInterface {
	x := s.term()
	var err error

	for s.matchLexeme("+") == nil || s.matchLexeme("-") == nil {
		tok := lexer.NewTokenWithString(s.cur_tok.Tag, s.lexer.Lexeme)
		s.move_forward()
		x, err = inter.NewArith(s.lexer.Line, tok, x, s.term())
		if err != nil {
			panic(err)
		}

	}

	return x
}

func (s *SimpleParser) term() inter.ExprInterface {
	x := s.factor()
	return x
}

func (s *SimpleParser) factor() inter.ExprInterface {
	var x inter.ExprInterface
	tok := lexer.NewTokenWithString(s.cur_tok.Tag, s.lexer.Lexeme)
	if s.matchTag(lexer.NUM) == nil {
		t := inter.NewType("int", lexer.BASIC, 4)
		x = inter.NewConstant(s.lexer.Line, tok, t)
	} else if s.matchTag(lexer.REAL) == nil {
		t := inter.NewType("float", lexer.BASIC, 8)
		x = inter.NewConstant(s.lexer.Line, tok, t)
	} else {
		sym := s.top.Get(s.lexer.Lexeme)
		if sym == nil {
			err_s := fmt.Sprintf("undefined variable with name: %s", s.lexer.Lexeme)
			err := errors.New(err_s)
			panic(err)
		}

		x = sym.id
	}

	s.move_forward()
	return x
}
