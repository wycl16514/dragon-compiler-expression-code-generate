package main

import (
	"lexer"

	"inter"
)

func main() {

	expr_type := inter.NewType("int", lexer.BASIC, 4)
	id_a := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "a"), expr_type)
	id_b := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "b"), expr_type)
	//a+b
	arith1, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), id_a, id_b)

	id_c := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "c"), expr_type)
	id_d := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "d"), expr_type)
	arith2, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), id_c, id_d)

	arith3, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.MINUS, "-"), arith1, arith2)

	//arith3.Reduce()

	id_e := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "e"), expr_type)
	//e = (a+b) - (b+c) -> c = t1 - t2
	set, _ := inter.NewSet(id_e, arith3)
	set.Gen()

	/*
	   	source := `{int x; float y ; float c; float d;
	              x = 1; y = 3.14;
	              c = x + y;
	              d = x + y + c;
	              }`
	   	my_lexer := lexer.NewLexer(source)
	   	parser := simple_parser.NewSimpleParser(my_lexer)
	   	parser.Parse()
	*/

}
