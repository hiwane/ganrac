package ganrac

type reduce_info struct {
	depth int
	q     []Level
	qn    int
	vars  *List
	varb  []bool
	eqns  *List
}

type simpler interface {
	simplBasic(neccon, sufcon Fof) Fof           // 手抜き簡単化
	simplComm() Fof                              // 共通部分の括りだし
	simplFctr(g *Ganrac) Fof                     // CA を要求
	simplReduce(g *Ganrac, inf *reduce_info) Fof // 等式制約による簡約化

	// symbolic-numeric simplification
	simplNum(g *Ganrac, true_region, false_region *NumRegion) (Fof, *NumRegion, *NumRegion)
	get_homo_cond(conds [][]int, c []int) [][]int
	homo_reconstruct(lv Level, lvs Levels, sgn int) Fof
}

func logSimplFof(c Fof, g *Ganrac, eyec string) {
	if false {
		return
	}
	// f, err := funcCAD(g, "fnCAD", []any{c})
	// if err != nil {
	// 	fmt.Printf("errrrr %v\n", err)
	// 	panic("stop")
	// }
	// c = f.(Fof)
	// @5->@6
	g.log(10, 1, "simplFof() step%s: %v\n", eyec, c)
}

func (g *Ganrac) simplFof(c Fof, neccon, sufcon Fof) Fof {
	g.log(3, 1, "simpli %v, nec=%v, suf=%v\n", c, neccon, sufcon)
	c = c.simplFctr(g)
	logSimplFof(c, g, "@1")
	c.normalize()
	logSimplFof(c, g, "@2")
	inf := NewReduceInfo()
	c = c.simplReduce(g, inf)
	logSimplFof(c, g, "@3")

	for {
		cold := c
		c = c.simplComm()
		logSimplFof(c, g, "@4")
		c = c.simplBasic(neccon, sufcon)
		logSimplFof(c, g, "@5")
		c, _, _ = c.simplNum(g, nil, nil)
		logSimplFof(c, g, "@6")
		if c.Equals(cold) {
			break
		}
	}

	g.log(3, 1, "simplo %v\n", c)
	return c
}
