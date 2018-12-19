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

func TestFollow() {
	var grammar = "S -> E\nE -> E + T\nE -> T\nT -> T * F\nT -> F\nF -> (E)\nF -> num"
	var t_set = []string{
		"num",
	}
	pros := parser.GetProdction(grammar, t_set) // 解析产生式
	pros = parser.RemoveRecursive(pros) // 移除左递归
	pros = parser.TakeCommonLeft(pros) // 提取左公因

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

	followSet := parser.Follow(pros, "T")

	fmt.Print("E FOLLOW -> ")

	for _, sym := range followSet {

		switch sym.SymType() {
		case parser.SYM_TYPE_TERMINAL:
			fmt.Printf(sym.Sym())

		case parser.SYM_TYPE_N_TERMINAL:
			fmt.Printf(sym.Sym())

		case parser.SYM_TYPE_NIL:
			fmt.Print("ε")
		}
		
	}

}

func TestLL1() {
	ll1 := &parser.LL1{}
	ll1.Init()
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

	// TestTakeCommonLeft()

	// TestFollow()
	TestLL1()
}