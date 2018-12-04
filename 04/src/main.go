package main

// import(
// 	"os"
// 	"bufio"
// 	"lex"
// 	"parser"
// )

import (
	"parser"
	"fmt"
)



func main() {
	// inputReader := bufio.NewReader(os.Stdin)
	// exp, _ := inputReader.ReadString('\n')
	// lex := &lex.Lex{}
	// lex.Init(exp)
	// tokens := lex.GetAllToken()

	// var parse parser.IParser = &parser.ReParser{}
	// parse.Init()
	// parse.SetTokens(tokens)
	// parse.Parse()

	var grammar = "S -> E\nE -> E + T\nE -> T\nT -> T * F\nT -> F\nF -> (E)\nF -> num"
	var t_set = []string{
		"num",
	}
	pros := parser.GetProdction(grammar, t_set)
	pros = parser.RemoveRecursive(pros)

	for _, value := range pros {
		fmt.Printf("%s -> ", value.Header())

		for _, sym := range value.Body() {

			switch sym.SymType() {
			case parser.SYM_TYPE_TERMINAL:
				fmt.Printf("%s%s%s", "T(", sym.Sym(), ")")

			case parser.SYM_TYPE_N_TERMINAL:
				fmt.Printf("%s%s%s", "NT(", sym.Sym(), ")")

			case parser.SYM_TYPE_NIL:
				fmt.Print("ε")
			}

			
		}

		fmt.Print("\n")

	}
}