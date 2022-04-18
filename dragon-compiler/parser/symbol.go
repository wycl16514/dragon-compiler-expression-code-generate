package simple_parser

import (
	"inter"
)

type Symbol struct {
	id        *inter.ID
	expr_type *inter.Type
}

func NewSymbol(id *inter.ID, expr_type *inter.Type) *Symbol {
	return &Symbol{
		id:        id,
		expr_type: expr_type,
	}
}
