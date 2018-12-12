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

func TestRemoveRecursive() {
	var grammar = "S -> E\nE -> E + T\nE -> T\nT -> T * F\nT -> F\nF -> (E)\nF -> num"
	grammar = "S -> Aa\nS -> b\nA -> Ac\nA -> Sd\nA -> @";
	var t_set = []string{
		"a",
		"b",
		"c",
		"d",
	}
	pros := parser.GetProdction(grammar, t_set)
	pros = parser.RemoveRecursive(pros)

	for _, value := range pros {
		fmt.Printf("%s -> ", value.Header())

		for _, sym := range value.Body() {

			switch sym.SymType() {
			case parser.SYM_TYPE_TERMINAL:
				fmt.Printf(sym.Sym())

			case parser.SYM_TYPE_N_TERMINAL:
				fmt.Printf(sym.Sym())

			case parser.SYM_TYPE_NIL:
				fmt.Print("ε")
			}

			
		}

		fmt.Print("\n")

	}
}

func TestTakeCommonLeft() {
	var grammar = "S -> E\nE -> aaaaF\nE -> aabbF\nE -> aabcF\nE -> aabdF\nE -> aacdF\nE -> abcdF\nE -> cbbF\nE -> cbcF\nF -> k"
	var t_set = []string{
		"a",
		"b",
		"c",
		"d",
		"k",
	}
	pros := parser.GetProdction(grammar, t_set)
	pros = parser.TakeCommonLeft(pros)

	for _, value := range pros {
		fmt.Printf("%s -> ", value.Header())

		for _, sym := range value.Body() {

			switch sym.SymType() {
			case parser.SYM_TYPE_TERMINAL:
				fmt.Printf(sym.Sym())

			case parser.SYM_TYPE_N_TERMINAL:
				fmt.Printf(sym.Sym())

			case parser.SYM_TYPE_NIL:
				fmt.Print("ε")
			}
			
		}

		fmt.Print("\n")

	}
}

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

	TestTakeCommonLeft()
}