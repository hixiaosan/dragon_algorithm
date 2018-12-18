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
				production.body = append(production.body, &Symbolic{sym_type: SYM_TYPE_NIL, sym: "@"})
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
		if len(v) > 1 { // 重复前缀大于 1个的 加入子节点
	
			child := &Node{Child: make([]*Node, 0), Pros:v, Sym: k}
			parent.Child = append(parent.Child, child)

			InitProductionTree(child, groupPros[k], idx + 1)
		}
	}

}



// 提取左部公因子
func TakeCommonLeft(pros [] *Production) ([]*Production) {
	result := make([]*Production, 0)
	multPro := GroupProduction(pros)
	
	
	for i := 0; i < len(multPro);  {
		pro := multPro[i]
		var tree Tree
		tree.Root = &Node{Child: make([]*Node, 0)}

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

// 提取First集合
func First(cfg []*Production, sym *Symbolic) map[string] *Symbolic {
	result := make(map[string] *Symbolic)

	// 规则一 如果符号是一个终结符号，那么他的FIRST集合就是它自身
	if sym.SymType() == SYM_TYPE_TERMINAL || sym.SymType() == SYM_TYPE_NIL {
		result[sym.Sym()] = sym
		return result
	}

	// 规则二 如果一个符号是一个非终结符号
	// (1) A -> XYZ 如果 X 可以推导出nil 那么就去查看Y是否可以推导出nil
	//              如果 Y 推导不出nil，那么把Y的First集合加入到A的First集合
	//				如果 Y 不能推导出nil，那么继续推导 Z 是否可以推导出nil,依次类推
	// (2) A -> XYZ 如果XYZ 都可以推导出 nil, 那么说明A这个产生式有可能就是nil， 这个时候我们就把nil加入到FIRST(A)中

	for _, production := range cfg { 

		if production.header == sym.Sym() {

			nilCount := 0
			for _, rightSymbolic := range production.body { // 对于一个产生式
				ret := First(cfg, rightSymbolic) // 获取这个产生式体的First集合
				hasNil := false
				for k, v := range ret {
					if v.SymType() == SYM_TYPE_NIL { // 如果推导出nil, 标识当前产生式体的符号可以推导出nil
						hasNil = true
					} else {
						result[k] = v
					}
				}

				if false == hasNil {  // 当前符号不能推导出nil, 那么这个产生式的FIRST就计算结束了，开始计算下一个产生式
					break
				} 

				// 当前符号可以推导出nil，那么开始推导下一个符号
				nilCount++

				if nilCount == len(production.body) { // 如果产生式体都可以推导出nil，那么这个产生式就可以推导出nil
					result["@"] = &Symbolic{sym: "@", sym_type: SYM_TYPE_NIL}
				}

			}



		}

	}

	return result
}

// 提取FOLLOW集合
func Follow(cfg []*Production, sym string) [] *Symbolic {
	fmt.Printf("Follow ------> %s\n", sym)
	result := make([] *Symbolic, 0)

	// 一个文法符号的FOLLOW集就是 可能出现在这个文法符号后面的终结符
	// 比如 S->ABaD, 那么FOLLOW(B)的值就是a。 
	//		            FOLLOW(A)的值包含了FIRST(B)中除了ε以外的所有终结符,如果First(B)包含空的话。说明跟在B后面的终结符号就可以跟在A后面，这时我们要把FOLLOW(B)的值也添加到FOLLOW(A)中
	//                  因为D是文法符号S的最右符号，那么所有跟在S后面的符号必定跟在D后面。所以FOLLOW(S)所有的值都在FOLLOW(D)中
	// 					以下是书中的总结

	// 不断应用下面这两个规则，直到没有新的FOLLOW集 被加入
	// 规则一: FOLLOW(S)中加入$, S是文法开始符号
	// 规则二: A->CBY FOLLOW(B) 就是FIRST(Y)
	// 规则三: A->CB 或者 A->CBZ(Z可以推导出ε) 所有FOLLOW(A)的符号都在FOLLOW(B), 
	//        

	if sym == "S" { // 如果是文法的开始符号
		result = append(result, &Symbolic{sym: "$", sym_type: SYM_TYPE_TERMINAL})
	}

	for _, pro := range cfg {
		for idx, p_r_sym := range pro.body {
			if p_r_sym.Sym() == sym { // 寻找到这个符号
				if idx + 1 == len(pro.body) { // 是文法最右部的符号

					if pro.header == p_r_sym.Sym() {
						continue
					}

					ret := Follow(cfg, pro.header) // 获取产生式头的FOLLOW集合
					result = append(result, ret...)
					continue
				}

				firstSet := First(cfg, pro.body[idx + 1]) // 获取下一个First集合

				hasNil := false
				for _, first := range firstSet {
					fmt.Printf("first -> %s\n", first.Sym())
					if first.SymType() == SYM_TYPE_NIL {
						hasNil = true
						continue
					}

					result = append(result, first) // 添加Follow集合
				}

				if hasNil { // 如果下一个符号包含空串 那么要在获取下一个的符号的Follow集合
					nextFollow := Follow(cfg, pro.body[idx + 1].Sym())
					result = append(result, nextFollow...)
					continue
				}
			}
		}
		
	}

	unique := make(map[string]*Symbolic)

	for _, sym := range result {
		if _, ok := unique[sym.Sym()]; !ok {
			unique[sym.Sym()] = sym
		}
	}

	result = make([]*Symbolic, 0)
	for _, sym := range unique {
		result = append(result, sym)
	}

	return result
}