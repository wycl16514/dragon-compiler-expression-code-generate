module dragon-compiler

go 1.17

require (
	inter v0.0.0-00010101000000-000000000000 // indirect
	simple_parser v0.0.0-00010101000000-000000000000
)

require lexer v0.0.0-00010101000000-000000000000

replace lexer => ./lexer

replace inter => ./inter

replace simple_parser => ./parser
