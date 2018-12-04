package lex

import (
	"fmt"
	"unsafe"
	"os"
)

const (
	TOKEN_NUM = 0
	TOKEN_ID = 1
	TOKEN_OPER = 2
	TOKEN_OPEN_EXP = 3
	TOKEN_CLOSE_EXP = 4
	TOKEN_END = 5
)

type Lex struct {
	idx int
	look_idx int
	str string
	data []byte
}



type TOKEN struct {
	TokenType int
	TokenAttr int
}


func (lex *Lex)Init(data string) {
	lex.str = data
	lex.data = []byte(lex.str)
	lex.data = append(lex.data, '$')
	fmt.Print(string(lex.data))
	lex.idx = 0
	lex.look_idx = 0
}

func (lex *Lex)get_char() byte {
	ch := lex.data[lex.look_idx]
	lex.look_idx++

	return ch
}

func (lex *Lex)get_token_num() (*TOKEN) {

	for {
		ch := lex.get_char()

		if ch < '0' || ch > '9' { // NUM TOKEN END
			val := lex.data[lex.idx : lex.look_idx]
			num_attr := string(val)
			var p uintptr = uintptr(unsafe.Pointer(&num_attr))
			lex.look_idx--
			lex.idx = lex.look_idx
			return &TOKEN{TokenType:TOKEN_NUM, TokenAttr:(int)(p)}
		}
	}
	
}

func (lex *Lex)filter_wp() {
	for {
		ch := lex.get_char()

		if ch != ' ' && ch != '\r' && ch != '\n' && ch != '\t' {
			lex.look_idx--
			lex.idx = lex.look_idx
			return 
		}
	}
	

}

// 对算数表达式的词法分析
func (lex *Lex)get_token() (*TOKEN) {
	
	ch := lex.get_char()

	if ch == ' ' || ch == '\r' || ch == '\n' || ch == '\t' {
		lex.filter_wp()
		ch = lex.get_char()
	}

	
	if ch == '$' {
		return &TOKEN{TokenType: TOKEN_END}
	}

	if ch >= '0' && ch <= '9' {
		return lex.get_token_num()
	}

	switch ch {
	case '+':
		fallthrough
	case '-':
		fallthrough
	case '*':
		fallthrough
	case '/':
		return &TOKEN{TokenType: TOKEN_OPER, TokenAttr: (int)(ch)}

	case '(':
		return &TOKEN{TokenType: TOKEN_OPEN_EXP}

	case ')':
		return &TOKEN{TokenType: TOKEN_CLOSE_EXP}
	}

	fmt.Print("词法分析错误")
	os.Exit(0)

	return nil
}

func (lex *Lex)GetAllToken() ([]*TOKEN) {
	tokens := make([]*TOKEN, 0, 0)
	for {
		token := lex.get_token()
		tokens = append(tokens, token)

		if (token.TokenType == TOKEN_END) {
			fmt.Print("词法分析成功\n")
			return tokens
		}

		
	}
}
