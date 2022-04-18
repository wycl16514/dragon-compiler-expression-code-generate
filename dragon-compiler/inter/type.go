package inter

import (
	"lexer"
)

type Type struct {
	width  uint32    //用多少字节存储该类型
	tag    lexer.Tag //
	Lexeme string
}

func NewType(lexeme string, tag lexer.Tag, w uint32) *Type {
	return &Type{
		width:  w,
		tag:    tag,
		Lexeme: lexeme,
	}
}

func Numberic(p *Type) bool {
	//查看给定类型是否属于数值类
	numberic := false
	switch p.Lexeme {
	case "int":
		numberic = true
	case "float":
		numberic = true
	case "char":
		numberic = true
	}

	return numberic
}

func MaxType(p1 *Type, p2 *Type) *Type {
	/*比较类型提升，例如p1是int，p2是float, 那么就提升为float
	类型提升必须对数值类型才有效
	*/

	if Numberic(p1) == false && Numberic(p2) == false {
		return nil
	}
	//如果两者有其一是float类型，那么就提升为float，要不然就是int
	if p1.Lexeme == "float" || p2.Lexeme == "float" {
		return NewType("float", lexer.BASIC, 8)
	} else if p1.Lexeme == "int" || p2.Lexeme == "int" {
		return NewType("int", lexer.BASIC, 4)
	}

	return NewType("char", lexer.BASIC, 1)
}
