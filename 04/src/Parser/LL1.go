package parser
import (
	"fmt"
)

// LL(1) 语法分析
type LL1 struct {
	p_table map[string]map[string]*Production // 预测分析表
}

// 传入的参数是文法产生式
func (parser *LL1)init_table(cfg []*Production) {
	parser.p_table = make(map[string]map[string]*Production)
	for _, pro := range cfg {
		if _, ok := parser.p_table[pro.header]; !ok {
			parser.p_table[pro.header] = make(map[string]*Production)
		}

		hasNil := false

		for _, pro_sym := range pro.body {
			first := First(cfg, pro_sym)

			for _, sym := range first {

				if sym.SymType() == SYM_TYPE_NIL {
					hasNil = true
					continue
				}

				parser.p_table[pro.header][sym.Sym()] = pro
			}

			if false == hasNil {
				break
			}
		}
		

		// 如果这个产生式可以推导出空串，那么我们就获取他的FOLLOW集
		if hasNil {
			follow := Follow(cfg, pro.header)
			for _, sym := range follow {
				parser.p_table[pro.header][sym.Sym()] = pro
			}
		}
	}

	fmt.Print("预测分析表\n")
	// 打印预测分析表
	for k, _map := range parser.p_table {
		fmt.Printf("%s -> M", k)
		fmt.Printf("\n")

		for _k, value := range _map {
			fmt.Printf("%s -> (", _k)

	
			fmt.Printf("%s -> ", value.Header())

			for _, sym := range value.Body() {

				switch sym.SymType() {
				case SYM_TYPE_TERMINAL:
					fmt.Printf(sym.Sym())

				case SYM_TYPE_N_TERMINAL:
					fmt.Printf(sym.Sym())

				case SYM_TYPE_NIL:
					fmt.Print("ε")
				}
				
			}

			fmt.Print(")  \n")
		}

		fmt.Print("\n\n")
	}

}

func (parser *LL1)Init() {
	var grammar = "S -> E\n" +
				  "E -> E + T\n" +
				  "E -> T\n" + 
				  "T -> T * F\n" + 
				  "T -> F\n" + 
				  "F -> (E)\n " + 
				  "F -> num"

	var t_set = []string{
		"num",
	}
	pros := GetProdction(grammar, t_set) // 解析文法
	pros = RemoveRecursive(pros)	// 消除左递归
	pros = TakeCommonLeft(pros)		// 提取左公因

	parser.init_table(pros)
	
}

func (parser *LL1)Parse() {

}