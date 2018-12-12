package parser
import (
	"fmt"
)

// 符号类型
const (
	SYM_TYPE_TERMINAL = 0
	SYM_TYPE_N_TERMINAL = 1
	SYM_TYPE_NIL = 2 // 空串
)

// 文法符号
type Symbolic struct {
	sym_type int // 符号类型
	sym      string // 符号名称
}

func (s *Symbolic)SymType() int {
	return s.sym_type
}

func (s *Symbolic)Sym() string {
	return s.sym
}

// 产生式
type Production struct {
	header string // 产生式头
	body   []*Symbolic // 产生式体
}

func (p *Production)Header() string {
	return p.header
}

func (p *Production)Body() []*Symbolic {
	return p.body
}

const (
	LEX_GRAMMAR_S_HEAD = 0
	LEX_GRAMMAR_S_BODY = 1
)

const (
	LEX_NT_START = 0
	LEX_NT_END = 1
)

// 解析文法 获取产生式
func GetProdction(grammar string, tset []string) ([]*Production) {

	var NTSet map[string]bool // 非终结符集合
	var NT string // 终结符
	var symbolic string // 符号
	var production *Production // 产生式
	var ProductList []*Production // 产生式集合

	ProductList = make([]*Production, 0, 0)
	NTSet = make(map[string]bool)

	data := []byte(grammar)

	status := LEX_NT_START

	// 获取所有的非终结符
	for i := 0; i < len(data); i++ {

		if data[i] == ' ' {
			continue
		}

		if data[i] == '-'  {
			status = LEX_NT_END
			NTSet[NT] = true
			NT = ""
			continue
		}

		if data[i] == '\r' || data[i] == '\n' {
			status = LEX_NT_START
			continue
		}

		if status == LEX_NT_START {
			NT += string(data[i]);
			continue;
		}

	}


	status = LEX_GRAMMAR_S_HEAD

	for i := 0; i < len(data); i++ {
		if data[i] == ' ' || data[i] == '>' {
			continue
		}

		if data[i] == '-'  {
			production = &Production{}
			production.header = symbolic
			ProductList = append(ProductList, production)
			status = LEX_GRAMMAR_S_BODY
			symbolic = ""
			continue
		}

		if data[i] == '\r' || data[i] == '\n' {
			status = LEX_GRAMMAR_S_HEAD
			continue
		}

		switch status {
		case LEX_GRAMMAR_S_HEAD:
			symbolic += string(data[i])

		case LEX_GRAMMAR_S_BODY: // 获取产生式体
			var S_B = i
			var match_len = 0
			var NTS string = "" // 非中介符号
			// 获取一个符号
			for key, _ := range NTSet {
				var tmp_len = 0
				for nti := 0; nti < len(key); nti++ {

					if key[nti] == data[S_B + nti] {
						tmp_len++
					} else {
						tmp_len = 0
						break
					}

				}

				if tmp_len > match_len {
					match_len = tmp_len
					NTS = key
				}

			}

			if NTS != "" {
				i = i + match_len - 1
				production.body = append(production.body, &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: NTS})
			} else if (data[i] == '@') {
				production.body = append(production.body, &Symbolic{sym_type: SYM_TYPE_NIL})
			} else {

				TS := ""

				for _, key := range tset {
					var tmp_len = 0
					for nti := 0; nti < len(key); nti++ {

						if key[nti] == data[S_B + nti] {
							tmp_len++
						} else {
							tmp_len = 0
							break
						}

					}

					if tmp_len > match_len {
						match_len = tmp_len
						TS = key
					}

				}

				if TS != "" {
					i = i + match_len - 1
					production.body = append(production.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: TS})
				} else {
					production.body = append(production.body, &Symbolic{sym_type: SYM_TYPE_TERMINAL, sym: string(data[i])})
				}

				
			}
		}
	}

	return ProductList
}

type MultProduction struct {
	header string
	pros []*Production
}

// 简化产生式
func SimpleProduction(pros []*Production) []*Production {
	result := make([]*Production, 0)
	
	return result
}

// 组合产生式，组合具有相同左部的产生式
func GroupProduction(pros []*Production) []*MultProduction {
	group := make(map[string][]*Production)
	multPro := make([]*MultProduction, 0)
	
	fmt.Println(len(pros))
	for i := 0; i < len(pros); i++ {
		if _, ok := group[pros[i].header]; !ok {
			fmt.Println(pros[i])
			group[pros[i].header] = make([]*Production, 0)
		}

		group[pros[i].header] = append(group[pros[i].header], pros[i])
	}

	for k, v := range group {
		mp := &MultProduction{header: k, pros: v}
		multPro = append(multPro, mp)
	}
	return multPro
}

// 左递归移除算法
func RemoveRecursive(pros []*Production) []*Production {

	result := make([]*Production, 0)
	multPro := GroupProduction(pros)

	for i := 0; i < len(multPro); i++ {
		for j := 0; j < i; j++ {
			newPros := make([]*Production, 0)

			for t := 0; t < len(multPro[i].pros); t++ {
				// 第一个符号不是终结符
				if multPro[i].pros[t].body[0].SymType() != SYM_TYPE_N_TERMINAL {
					newPros = append(newPros, multPro[i].pros[t])
					continue
				}

				// 立即左递归
				if multPro[i].pros[t].body[0].Sym() == multPro[i].header {
					newPros = append(newPros, multPro[i].pros[t])
					continue
				}

				if multPro[i].pros[t].body[0].Sym() == multPro[j].header {

					for c := 0; c < len(multPro[j].pros); c++ {
						pros := &Production{header: multPro[j].header}
						pros.body = make([]*Symbolic, 0)
						pros.body = append(pros.body, multPro[j].pros[c].body...)
						pros.body = append(pros.body, multPro[i].pros[t].body[1:]...)

						newPros = append(newPros, pros)
					}

				} else { // 必须是 Ai -> Ajγ 不匹配的话加入原本的产生式

					newPros = append(newPros, multPro[i].pros[t])
				}
			}

			multPro[i].pros = newPros // 设置 Ai 为替换后的产生式

		}

		var isLeftRecursive bool = false

		// 消除直接左递归
		for t := 0; t < len(multPro[i].pros); t++ {
			
			// 存在直接左递归
			if multPro[i].pros[t].body[0].Sym() == multPro[i].header {
				isLeftRecursive = true
			}

		}

		// 没有直接左递归 拷贝产生式
		if isLeftRecursive == false {
			result = append(result, multPro[i].pros...)
			continue
		}

		for t := 0; t < len(multPro[i].pros); t++ {
			
			// 左递归产生式
			if multPro[i].pros[t].body[0].Sym() == multPro[i].header {
				pro := &Production{header: multPro[i].header + "`"}
				pro.body = make([]*Symbolic, 0)
				pro.body = append(pro.body, multPro[i].pros[t].body[1:]...)
				pro.body = append(pro.body, &Symbolic{sym: multPro[i].header + "`", sym_type: SYM_TYPE_N_TERMINAL})
				result = append(result, pro)
			} else {
				// 非递归产生式

				// 空产生式 A->ε 变成 A->A`
				if multPro[i].pros[t].body[0].SymType() == SYM_TYPE_NIL {
					pro := &Production{header: multPro[i].header}
					pro.body = make([]*Symbolic, 0)
					pro.body = append(pro.body, &Symbolic{sym: multPro[i].header + "`", sym_type: SYM_TYPE_N_TERMINAL})
					result = append(result, pro)
				} else {
					// 非空产生式 A->ab 变成 A->abA`
					pro := &Production{header: multPro[i].header}
					pro.body = make([]*Symbolic, 0)
					pro.body = append(pro.body, multPro[i].pros[t].body...)
					pro.body = append(pro.body, &Symbolic{sym: multPro[i].header + "`", sym_type: SYM_TYPE_N_TERMINAL})
					result = append(result, pro)
				}

			}

		}

		// 加入一个空产生式
		pro := &Production{header: multPro[i].header + "`"}
		pro.body = make([]*Symbolic, 0)
		pro.body = append(pro.body, &Symbolic{sym: "", sym_type: SYM_TYPE_NIL})
		result = append(result, pro)

	}

	return result

	//return SimpleProduction(result)
}

type Tree struct {
	Root *Node
}

type Node struct {
	Child []*Node  // 子节点
	Pros []*Production	   // 产生式索引
	Sym  Symbolic // 当前符号
}

// 生成分析树
func InitProductionTree(parent *Node, pros []*Production, idx int) {

	fmt.Println("InitProductionTree")
	groupPros := make(map[Symbolic][]*Production)

	for i := 0; i < len(pros); i++ {
		if len(pros[i].body) > idx {
			sym := *pros[i].body[idx]
			if _, ok := groupPros[sym]; !ok {
				groupPros[sym] = make([]*Production, 0);
			}

			groupPros[sym] = append(groupPros[sym], pros[i])
		}
	}

	for k, v := range groupPros {
		fmt.Println(k)
		if len(v) > 1 { // 重复前缀大于 1个的 加入子节点
	
			child := &Node{Child: make([]*Node, 0), Pros:v, Sym: k}
			parent.Child = append(parent.Child, child)

			InitProductionTree(child, groupPros[k], idx + 1)
		}
	}
	fmt.Println(len(parent.Child))


}



// 提取左部公因子
func TakeCommonLeft(pros [] *Production) ([]*Production) {
	result := make([]*Production, 0)
	multPro := GroupProduction(pros)
	
	
	for i := 0; i < len(multPro);  {
		pro := multPro[i]
		var tree Tree
		tree.Root = &Node{Child: make([]*Node, 0)}

		fmt.Println("LOOP: ---------------------", i)
		InitProductionTree(tree.Root, pro.pros, 0)
		if len(tree.Root.Child) == 0 { // 没有公共前缀
			i++
			
			continue
		}

		// 有公共前缀的处理
		root := tree.Root
		last := root
		var deep int  = 0 // 层次
		preSymbolic := make([]Symbolic, 0)
		for  {
			if len(root.Child) == 1 { // 前缀全部相同, 进入下一层
				last = root
				root = root.Child[0]
				preSymbolic = append(preSymbolic, root.Sym) // 前缀符号
				deep++
				continue
			}

			if len(root.Child) > 1 || (len(root.Child) == 0 && deep != 0) {
				
				// 前缀不是全部相同

				nmPros := &MultProduction{header: pro.header + "`", pros: make([]*Production, 0)} // 新的产生式集合
				fmt.Println("----------------------,", len(root.Child), deep)
				
				if deep != 0 { // 不是第一层
					newProc := &Production{header: pro.header}

					for _, sym := range preSymbolic {
						newProc.body = append(newProc.body, &Symbolic{sym_type: sym.sym_type, sym: sym.sym})
					}

					newProc.body = append(newProc.body, &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: newProc.header + "`"})

					newPros := make([]*Production, 0)

					root = last // 回到上一层
					newPros = append(newPros, newProc)

					for _, pro := range multPro[i].pros {
						NTProc := &Production{header: pro.header + "`", body: make([]*Symbolic, 0)}
						fmt.Println(pro.header, pro.body[0])
						NTProc.body = append(NTProc.body, pro.body[deep:]...)
						nmPros.pros = append(nmPros.pros, NTProc)
					}

					multPro[i].pros = newPros
				} else { // 第一层

					for _, c := range root.Child {
						preSymbolic = preSymbolic[0:0]
						preSymbolic = append(preSymbolic, c.Sym) // 前缀符号

						newProc := &Production{header: pro.header}

						for _, sym := range preSymbolic {
							newProc.body = append(newProc.body, &Symbolic{sym_type: sym.sym_type, sym: sym.sym})
						}

						newProc.body = append(newProc.body, &Symbolic{sym_type: SYM_TYPE_N_TERMINAL, sym: newProc.header + "`"})

						newPros := make([]*Production, 0)

						for _, org_pro := range multPro[i].pros {
							find := false
							for _, del := range c.Pros {
		
								if del == org_pro {
									find = true
									NTProc := &Production{header: pro.header + "`", body: make([]*Symbolic, 0)}
									NTProc.body = append(NTProc.body, del.body[1:]...)

									// 防止符号重复
									exists := false
									f:for _, pro := range nmPros.pros {
										if len(pro.body) != len(NTProc.body) {
											continue
										}

										for _t := 0; _t < len(pro.body); _t++ {
											if *pro.body[_t] != *NTProc.body[_t] {
												break 
											}

											if _t + 1 == len(pro.body) {
												exists = true
												break f
											}
										}
									}

									if false == exists {
										nmPros.pros = append(nmPros.pros, NTProc)
									}
									break
								}
							}

							if find == false {
								newPros = append(newPros, org_pro)
							}
						}

						newPros = append(newPros, newProc)
						multPro[i].pros = newPros
					}
					
				}

				multPro = append(multPro, nmPros)
				break
			}
			

		}
	}


	for i := 0; i < len(multPro);  i++ {
		result = append(result, multPro[i].pros...)
	}


	return result
}