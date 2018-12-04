package parser

import (
	"fmt"
	"os"
	"lex"
)

type N_T_Process func (parser *ReParser)


// 递归下降的语法分析
type ReParser struct {
	idx int
	tokens []*lex.TOKEN	  // 输入
	cfg []*Production  // 文法
	process map[string]N_T_Process
	start string
}

// 初始化文法
func (parser *ReParser)Init() {
	E := &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: "E"}
	T := &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: "T"}
	E_ := &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: "E`"}
	F := &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: "F"}
	T_ := &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: "T`"}
	NIL := &Symbolic{sym_type: SYM_TYPE_NIL, sym: ""}

	P_S := &Production{header: "S", body: make([]*Symbolic, 0, 0)}
	P_S.body = append(P_S.body, E)

	P_E := &Production{header: "E", body: make([]*Symbolic, 0, 0)}
	P_E.body = append(P_E.body, T)
	P_E.body = append(P_E.body, E_)

	P_E_1 := &Production{header: "E`", body: make([]*Symbolic, 0, 0)}
	P_E_1.body = append(P_E_1.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: "+"})
	P_E_1.body = append(P_E_1.body, T)
	P_E_1.body = append(P_E_1.body, E_)

	P_E_2 := &Production{header: "E`", body: make([]*Symbolic, 0, 0)}
	P_E_2.body = append(P_E_2.body, NIL)

	P_T := &Production{header: "T", body: make([]*Symbolic, 0, 0)}
	P_T.body = append(P_T.body, F)
	P_T.body = append(P_T.body, T_)

	P_T_1 := &Production{header: "T`", body: make([]*Symbolic, 0, 0)}
	P_T_1.body = append(P_T_1.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: "*"})
	P_T_1.body = append(P_T_1.body, F)
	P_T_1.body = append(P_T_1.body, T_)

	P_T_2 := &Production{header: "T`", body: make([]*Symbolic, 0, 0)}
	P_T_2.body = append(P_T_2.body, NIL)

	P_F1 := &Production{header: "F", body: make([]*Symbolic, 0, 0)}
	P_F1.body = append(P_F1.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: "("})
	P_F1.body = append(P_F1.body, E)
	P_F1.body = append(P_F1.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: ")"})

	P_F2 := &Production{header: "F", body: make([]*Symbolic, 0, 0)}
	P_F2.body = append(P_F2.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: "num"})

	parser.cfg = make([]*Production, 0, 0)
	parser.cfg = append(parser.cfg, P_S)
	parser.cfg = append(parser.cfg, P_E)
	parser.cfg = append(parser.cfg, P_E_1)
	parser.cfg = append(parser.cfg, P_E_2)
	parser.cfg = append(parser.cfg, P_T)
	parser.cfg = append(parser.cfg, P_T_1)
	parser.cfg = append(parser.cfg, P_T_2)
	parser.cfg = append(parser.cfg, P_F1)
	parser.cfg = append(parser.cfg, P_F2)

	parser.process = make(map[string]N_T_Process)
	parser.process["T"] = ParserT
	parser.process["E"] = ParserE
	parser.process["E`"] = ParserE_
	parser.process["T`"] = ParserT_
	parser.process["F"] = ParserF

	parser.start = "S"
}

func ParserT(parser *ReParser) {
	fmt.Println("ParserT")
	pros := parser.GetProducts("T")

	if len(pros) == 0 {
		fmt.Print("产生式T为空")
		return
	}

	// 当前的输入串

	for pi := 0; pi < len(pros); pi++ {

		for si := 0; si < len(pros[pi].body); si++ {

			sym := pros[pi].body[si]

			if sym.sym_type == SYM_TYPE_N_TERMINAL {
				parser.process[sym.sym](parser)
			} else {
				fmt.Print("产生式T错误")
				os.Exit(0)
			}
		}

	}
}

func ParserE_(parser *ReParser) {
	fmt.Println("ParserE_")
	pros := parser.GetProducts("E`")

	if len(pros) == 0 {
		fmt.Print("产生式T为空")
		return
	}

	// 当前的输入串
	err := false
	idx := parser.idx
	si := 0

	for pi := 0; pi < len(pros); pi++ {
		err = false
		for si = 0; si < len(pros[pi].body); si++ {

			sym := pros[pi].body[si]

			if sym.sym_type == SYM_TYPE_N_TERMINAL {
				parser.process[sym.sym](parser)
			} else if sym.sym_type == SYM_TYPE_TERMINAL {
				
				if sym.sym == "+" && 
				   parser.tokens[parser.idx].TokenType == lex.TOKEN_OPER && 
				   parser.tokens[parser.idx].TokenAttr == '+' {
					parser.idx++
				} else {
					parser.idx = idx
					err = true
					break // 尝试下一个产生式
				}

			} else { // 空串

			}
		}

		if false == err {
			return
		}
	}

	if err {
		fmt.Println("ParserE_ 语法分析失败")
		os.Exit(0)
	}
}

func ParserT_(parser *ReParser) {
	fmt.Println("ParserT_")
	pros := parser.GetProducts("T`")

	if len(pros) == 0 {
		fmt.Print("产生式T为空")
		return
	}

	// 当前的输入串
	idx := parser.idx
	si := 0
	err := false

	for pi := 0; pi < len(pros); pi++ {
		err = false
		for si = 0; si < len(pros[pi].body); si++ {

			sym := pros[pi].body[si]

			if sym.sym_type == SYM_TYPE_N_TERMINAL {
				parser.process[sym.sym](parser)
			} else if sym.sym_type == SYM_TYPE_TERMINAL {
				
				if sym.sym == "*" && 
				   parser.tokens[parser.idx].TokenType == lex.TOKEN_OPER && 
				   parser.tokens[parser.idx].TokenAttr == '*' {
					parser.idx++
				} else {
					parser.idx = idx
					err = true
					break // 尝试下一个产生式
				}

			} else { // 空串

			}
		}

		if false == err {
			return
		}

	}

	if err {
		fmt.Println("ParserT_ 语法分析失败")
		os.Exit(0)
	}

}

func ParserF(parser *ReParser) {
	fmt.Println("ParserF")
	pros := parser.GetProducts("F")

	if len(pros) == 0 {
		fmt.Print("产生式F为空")
		return
	}

	// 当前的输入串
	idx := parser.idx
	si := 0
	err := false

	for pi := 0; pi < len(pros); pi++ {
		err = false
		for si = 0; si < len(pros[pi].body); si++ {

			sym := pros[pi].body[si]

			if sym.sym_type == SYM_TYPE_N_TERMINAL {
				parser.process[sym.sym](parser)
			} else if sym.sym_type == SYM_TYPE_TERMINAL {

				if sym.sym == "(" && 
				   parser.tokens[parser.idx].TokenType == lex.TOKEN_OPEN_EXP {
					parser.idx++
				} else if  sym.sym == ")" && 
				   parser.tokens[parser.idx].TokenType == lex.TOKEN_CLOSE_EXP {
					parser.idx++
					fmt.Println("success")
				} else if  sym.sym == "num" && 
				   parser.tokens[parser.idx].TokenType == lex.TOKEN_NUM {
					parser.idx++
				} else {
					parser.idx = idx
					err = true
					break // 尝试下一个产生式
				}

			} else { // 空串

			}
		}

		

	}

	if false == err {
		return
	}

	if err {
		fmt.Println("ParseF 分析失败")
		os.Exit(0)
	}

}

func ParserE(parser *ReParser) {
	fmt.Println("ParserE")
	pros := parser.GetProducts("E")

	if len(pros) == 0 {
		fmt.Print("产生式E为空")
		return
	}

	// 当前的输入串

	for pi := 0; pi < len(pros); pi++ {

		for si := 0; si < len(pros[pi].body); si++ {

			sym := pros[pi].body[si]

			if sym.sym_type == SYM_TYPE_N_TERMINAL {
				parser.process[sym.sym](parser)
			} else {
				fmt.Print("产生式E错误")
				os.Exit(0)
			}
		}

	}
}

func (parser *ReParser)SetTokens(tokens []*lex.TOKEN) {
	parser.tokens = tokens
	parser.idx = 0
}

func (parser *ReParser)GetProducts(N_T string) ([]*Production) {
	
	products := make([]*Production, 0, 0)

	for i := 0; i < len(parser.cfg); i++ {

		if parser.cfg[i].header == N_T {
			products = append(products, parser.cfg[i])
		}

	}

	return products
}

func (parser *ReParser)Parse() {
	pros := parser.GetProducts("S")

	if 1 != len(pros) {
		fmt.Print("语法分析错误: 未找到文法开始符号")
		os.Exit(0)
	}

	start := pros[0]

	for i := 0; i < len(start.body); i++ {

		sym := start.body[i]

		if sym.sym_type == SYM_TYPE_N_TERMINAL {
			parser.process[sym.sym](parser)
		} else {
			fmt.Print("语法分析错误: Parser 终结符")
			os.Exit(0)
		}

	}

	fmt.Print("语法分析成功\n")


}
