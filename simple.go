package ganrac

import "fmt"

type reduce_info struct {
	depth int
	q     []Level
	qn    int
	vars  *List
	varb  []bool // quantified variable?
	eqns  *List  // list of *Poly
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

func logSimplFof(c, neccon, sufcon Fof, g *Ganrac, eyec string) {
	if true {
		return
	}
	if true {
		dic := NewDict()
		dic.Set("var", NewInt(1))
		f, err := funcCAD(g, "fnCAD", []any{c, dic})
		if err != nil {
			fmt.Printf("errrrr %v\n", err)
			panic("stop")
		}
		fmt.Printf("simplFof(!) step%s: %v: %v\n", eyec, f, c)
	} else {
		//fmt.Printf("simplFof() step%s: %v\n", eyec, c)
		g.log(1, 1, "simplFof() step%s: %v\n", eyec, c)
	}
}

func (g *Ganrac) simplFof(c Fof, neccon, sufcon Fof) Fof {
	g.log(3, 1, "simpli %v, nec=%v, suf=%v\n", c, neccon, sufcon)
	inf := NewReduceInfo(g, c, neccon, sufcon)
	c = c.simplReduce(g, inf) // GB の結果により重複因子が生まれる可能性がある
	logSimplFof(c, neccon, sufcon, g, "@1")

	c = c.simplFctr(g)
	logSimplFof(c, neccon, sufcon, g, "@2")
	c.normalize()
	logSimplFof(c, neccon, sufcon, g, "@3")

	for i := 0; i < 10; i++ {
		cold := c
		c = c.simplComm()
		logSimplFof(c, neccon, sufcon, g, "@4")
		c = c.normalize()
		c = c.simplBasic(neccon, sufcon)
		logSimplFof(c, neccon, sufcon, g, "@5")
		c, _, _ = c.simplNum(g, nil, nil) // fctr 済みで重複因子がないことを仮定
		logSimplFof(c, neccon, sufcon, g, "@6")
		if c.Equals(cold) {
			break
		}
	}

	g.log(3, 1, "simplo %v\n", c)
	return c
}
