package parser

import (
	"lex"
)



type IParser interface {
	SetTokens(tokens []*lex.TOKEN)
	Parse()
	Init()
}