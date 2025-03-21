package ganrac

import (
	"fmt"
)

type QeExample struct {
	Input  Fof
	Output Fof
	Ref    string
	DOI    string
}

type qeExTable struct {
	name string
	f    func() *QeExample
}

func ExampleNames() []string {
	ret := make([]string, len(qeExampleTable))
	for i, s := range qeExampleTable {
		ret[i] = s.name
	}
	return ret
}

var qeExampleTable []qeExTable = []qeExTable{
	{"adam1", exAdam1},
	{"adam2-1", exAdam2_1},
	{"adam2-2", exAdam2_2},
	{"adam3", exAdam3},
	{"candj", exCandJ},
	{"catasph", exCatastropheSurfaceSphere},
	{"constcd", exConstCoord},
	{"cycle3", exCyclic3},
	{"easy7", exEasy7},
	{"hong93", exHong93},
	{"imo13-1", exImo13_1_5},
	{"kyoto2010h2", exKyoto2010h2},
	{"lampinen", exLampinen},
	{"lass", exLass},
	{"makepdf", exMakePdf},
	{"makepd2", exMakePdf2},
	{"mooea", exMooea},
	{"neq1", exNeq1},
	{"neq2", exNeq2},
	{"neq3", exNeq3},
	{"neq4", exNeq4},
	{"neq5", exNeq5},
	{"neq6", exNeq6},
	{"pl01", exPL01},
	{"portfolio", exPortfolio},
	{"quad", exQuad},
	{"quart", exQuart},
	{"sdc2", exSDC2},
	{"sdc3", exSDC3},
	{"sdc4", exSDC4},
	{"root2", exRoot2},
	{"root3", exRoot3},
	{"root4", exRoot4},
	{"root4-2", exRoot4_2},
	{"root4-3", exRoot4_3},
	{"root6", exSato1},
	{"wlog1", exWlog1},
	{"wo1", exWO1},
	{"wo2", exWO2},
	{"wo3", exWO3},
	{"wo4", exWO4},
	{"xaxis", exXAxisEllipse},
	{"whitney", exWhitneyUmbrella},
	{"wilson", exWilson},
}

func GetExampleFof(name string) *QeExample {
	if name == "" {
		fmt.Printf("\nlabel\t\t# free\t# q \tdeg(f)\tdeg(q)\tatom\n")
		fmt.Printf("============\t======\t======\t=======\t======\t======\n")
		for _, t := range qeExampleTable {
			q := t.f()
			v := q.Input.maxVar()
			fdeg := 0
			qdeg := 0
			fnum := 0
			qnum := 0
			for i := Level(0); i <= Level(v); i++ {
				d := q.Input.Deg(i)
				if q.Input.hasFreeVar(i) {
					fnum++
					if d > fdeg {
						fdeg = d
					}
				} else if q.Input.hasVar(i) {
					qnum++
					if d > qdeg {
						qdeg = d
					}
				}
			}

			fmt.Printf("%-10s\t%4d\t%4d\t%4d\t%4d\t%4d\n", t.name, fnum, qnum, fdeg, qdeg, q.Input.numAtom())
		}
		return nil
	}

	for _, t := range qeExampleTable {
		if t.name == name {
			// fmt.Printf("%S\n", t.f().Input)
			return t.f()
		}
	}
	return nil
}

func exAdam1() *QeExample {
	q := new(QeExample)
	q.Output = trueObj
	q.Input = NewQuantifier(true, []Level{0, 1}, NewFmlOrs(
		NewAtom(NewPolyCoef(0, 0, 1), GE),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -49719, 0, 50000), 0, 50000), GE),
		NewAtom(NewPolyCoef(1,
			NewPolyCoef(0, 0, 720000, 720000, 480000, 240000, 96000, 32200, 9200, 2225, 450, 75, 10, 1),
			0,
			NewPolyCoef(0, 0, 0, 0, 0, -3000, 1200, 2100, 1000, 275, 50, 6),
			0,
			NewPolyCoef(0, 0, 0, 3000, -6000, -2250, 300, 350, 100, 15),
			0,
			NewPolyCoef(0, -200, 2000, -1900, -600, 150, 100, 20),
			0,
			NewPolyCoef(0, 225, -350, -25, 50, 15),
			0,
			NewPolyCoef(0, -25, 10, 6), 0, 1), LT)))
	q.Ref = "Adam W. Strzebonski. Cylindrical Algebraic Decomposition using validated numerics. 2006"
	q.DOI = "10.1016/j.jsc.2006.06.004"
	return q
}

func exAdam2_1() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{0, 1, 2}, NewFmlOrs(
		NewAtom(NewPolyCoef(0, 0, 1), LT),
		NewAtom(NewPolyCoef(1, 0, 1), LT),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 4), 0, 4), GE),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 4, -3, -4, 4), NewPolyCoef(0, -4, 2, 8, -8), NewPolyCoef(0, 5, -12, 8)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 2, -4), NewPolyCoef(0, 0, 4, -4, 8), NewPolyCoef(0, 2, -4, -8), NewPolyCoef(0, -4, 8)), NewPolyCoef(1, NewPolyCoef(0, 0, -4, 5), NewPolyCoef(0, 4, 2, -12), NewPolyCoef(0, -3, 8, 8), NewPolyCoef(0, -4, -8), 4)), LE),
		NewAtom(NewPolyCoef(2, 0, 0, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2, 0, 2), NewPolyCoef(0, 0, 2, 0, -4), NewPolyCoef(0, 0, -4, 6)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, -4, 4), NewPolyCoef(0, 0, -4), NewPolyCoef(0, 0, 4)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 2, -4), NewPolyCoef(0, -2, 0, 6), NewPolyCoef(0, 0, -4), 2)), LE)))

	q.Output = trueObj
	q.Ref = "Adam W. Strzebonski. Cylindrical Algebraic Decomposition using validated numerics. 2006"
	q.DOI = "10.1016/j.jsc.2006.06.004"
	return q
}
func exAdam2_2() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{0, 1, 2}, NewFmlOrs(
		NewAtom(NewPolyCoef(0, 0, 1), LT),
		NewAtom(NewPolyCoef(0, -1, 1), GT),
		NewAtom(NewPolyCoef(1, 0, 1), LT),
		NewAtom(NewPolyCoef(1, -1, 1), GT), NewFmlAnds(
			NewAtom(NewPolyCoef(2, 0, 0, 0, 0, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1, 0, 1), NewPolyCoef(0, 0, 2, 0, -4), NewPolyCoef(0, -1, 0, 6), NewPolyCoef(0, 0, -4), 1)), LE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 0, -1, 0, 1), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, -4, 4)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 2, -4), NewPolyCoef(0, 0, -4, 6)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, -4), NewPolyCoef(0, 0, 4)), NewPolyCoef(1, 0, 0, -1, 0, 1)), LE), NewFmlOrs(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 4, -3, -4, 4), NewPolyCoef(0, -4, 2, 8, -8), NewPolyCoef(0, 5, -12, 8)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 2, -4), NewPolyCoef(0, 0, 4, -4, 8), NewPolyCoef(0, 2, -4, -8), NewPolyCoef(0, -4, 8)), NewPolyCoef(1, NewPolyCoef(0, 0, -4, 5), NewPolyCoef(0, 4, 2, -12), NewPolyCoef(0, -3, 8, 8), NewPolyCoef(0, -4, -8), 4)), LE),
				NewAtom(NewPolyCoef(2, 0, 0, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2, 0, 2), NewPolyCoef(0, 0, 2, 0, -4), NewPolyCoef(0, 0, -4, 6)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, -4, 4), NewPolyCoef(0, 0, -4), NewPolyCoef(0, 0, 4)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 2, -4), NewPolyCoef(0, -2, 0, 6), NewPolyCoef(0, 0, -4), 2)), LE)))))

	q.Output = trueObj
	q.Ref = "Adam W. Strzebonski. Cylindrical Algebraic Decomposition using validated numerics. 2006"
	q.DOI = "10.1016/j.jsc.2006.06.004"
	return q
}
func exAdam3() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2, 3}, NewQuantifier(true, []Level{1}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, -1, 1), GT),
		NewAtom(NewPolyCoef(3, 0, 1), GT),
		NewAtom(NewPolyCoef(0, 0, 1), GT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, -1), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, -2)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 1, -1), 0, 1)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 2)), NewPolyCoef(1, NewPolyCoef(0, 1, -1), 0, 1)), LE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, -1), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 2)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 1, -1), 0, 1)), NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 2)), NewPolyCoef(1, NewPolyCoef(0, 1, -1), 0, 1)), LE))))
	q.Input = NewQuantifier(false, []Level{1, 2}, NewQuantifier(true, []Level{3}, NewFmlAnds(
		NewAtom(NewPolyCoef(1, -1, 1), GT),
		NewAtom(NewPolyCoef(2, 0, 1), GT),
		NewAtom(NewPolyCoef(0, 0, 1), GT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(0, 1, -1)), 0, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), NewPolyCoef(0, 0, -2), NewPolyCoef(0, 1, -1)), NewPolyCoef(0, 0, 2), 1), 0, NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, 1)), LE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(0, 1, -1)), 0, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), NewPolyCoef(0, 0, 2), NewPolyCoef(0, 1, -1)), NewPolyCoef(0, 0, 2), 1), 0, NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, 1)), LE))))

	q.Output = NewAtom(NewPolyCoef(0, -4, 1), GT)
	q.Ref = "Adam W. Strzebonski. Cylindrical Algebraic Decomposition using validated numerics. 2006"
	q.DOI = "10.1016/j.jsc.2006.06.004"
	return q
}

func exCandJ() *QeExample {
	// ex([z], z>0 && z-1<0 && y>0 && 2*x >= 1 && (3*y^2+3*x^2-2*x)*z+-y^2-x^2<0 && (3*y^2+3*x^2-4*x+1)*z+-2*y^2+-2*x^2+2*x>0)

	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, 0, 1), GT),
		NewAtom(NewPolyCoef(2, -1, 1), LT),
		NewAtom(NewPolyCoef(1, 0, 1), GT),
		NewAtom(NewPolyCoef(0, -1, 2), GE),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 0, -1), NewPolyCoef(1, NewPolyCoef(0, 0, -2, 3), 0, 3)), LT),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 2, -2), 0, -2), NewPolyCoef(1, NewPolyCoef(0, 1, -4, 3), 0, 3)), GT)))

	q.Output = q.Input

	return q
}

func exCatastropheSurfaceSphere() *QeExample {
	// ex([x,y,z], z^2+y^2+x^2-1==0 && z^3+x*z+y==0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{0, 1, 2}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), 0, 1), EQ),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 0, 1), EQ)))
	q.Output = trueObj
	q.Ref = "Scott MaCallum. An Improved Projection Operation for Cylindrical Algebraic Decomposition"
	return q
}

func exConstCoord() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{3, 4}, NewFmlOrs(
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, -1, 1), 1), LE),
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, -1)), 0, NewPolyCoef(1, 0, 0, 1)), NewPolyCoef(1, 0, 0, 1)), NewPolyCoef(0, 0, 0, 1)), NE),
		NewAtom(NewPolyCoef(3, 0, 1), LT),
		NewAtom(NewPolyCoef(4, 0, 1), LT)))
	// (A w)(A a) [ a+w <= 1 \/ x^2 a + y^2 w + y^2 z^2 - x^2 y^2 /= 0 \/ w < 0 \/ a < 0].
	// y /= 0 /\ y^2 z^2 - x^2 y^2 + x^2 >= 0 /\ z^2 - x^2 + 1 >= 0 /\ [ x /= 0 \/ y^2 z^2 - x^2 y^2 + x^2 > 0 ]
	// y != 0 && y^2*z^2 - x^2*y^2 + x^2 >= 0 && z^2 - x^2 + 1 >= 0 && ( x != 0 || y^2*z^2 - x^2*y^2 + x^2 > 0 )
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(1, 0, 1), NE),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 1), 0, NewPolyCoef(0, 0, 0, -1)), 0, NewPolyCoef(1, 0, 0, 1)), GE),
		NewAtom(NewPolyCoef(2, NewPolyCoef(0, 1, 0, -1), 0, 1), GE),
		NewFmlOrs(
			NewAtom(NewPolyCoef(0, 0, 1), NE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 1), 0, NewPolyCoef(0, 0, 0, -1)), 0, NewPolyCoef(1, 0, 0, 1)), GT)))
	q.Ref = "original."
	return q
}

func exCyclic3() *QeExample {
	q := new(QeExample)
	q.Output = falseObj
	// ex([y,z], z+y+x==0 && (y+x)*z+x*y==0 && x*y*z-1==0)
	q.Input = NewQuantifier(false, []Level{1, 2}, NewFmlAnds(NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), 1), EQ), NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1)), NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1)), EQ), NewAtom(NewPolyCoef(2, -1, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1))), EQ)))
	return q
}

func exEasy7() *QeExample {
	// z-2*y-x==0 && z^2-125==0 && z>0 && y^2+x^2-25<=0 && y^2+(-x-5)*y-2*x^2-20*x-50<=0
	// sqrt(125) = 2*y+x && y^2+x^2 <= 25 && y^2+(-x-5)*y-2*x^2-20*x-50<=0
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), -2), 1), EQ),
		NewAtom(NewPolyCoef(2, -125, 0, 1), EQ),
		NewAtom(NewPolyCoef(2, 0, 1), GT),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -25, 0, 1), 0, 1), LE),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -50, -20, -2), NewPolyCoef(0, -5, -1), 1), LE)))
	q.Output = NewFmlAnds( // sncad x > 0 and y^2+x^2 == 25 and -4*y^2-4*x*y-x^2+125 == 0
		NewAtom(NewPolyCoef(0, 0, 1), GT),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -25, 0, 1), 0, 1), EQ),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -125, 0, 1),
			NewPolyCoef(0, 0, 4), 4), EQ))
	// redlog->qepcad: x > 0 /\ x^2 - 5 = 0 /\ 2 y + x > 0 /\ y^2 + x^2 - 25 = 0

	// alpha: x^2+y^2<=25 && x^2+4*x*y+4*y^2 == 125 && x+2*y > 0 && 2*x^2+x*y+20*x-y^2+5*y+50 >= 0
	q.Ref = "syn_pdq error"

	return q
}

func exKyoto2010h2() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{0, 1}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -2), -1), 1), EQ),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -9, 4), 1), LE),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -4, 1), 2), GE),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -6, -2), 3), LE)))
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(2, -6, 1), LE),
		NewAtom(NewPolyCoef(2, -2, 1), GE))
	q.Ref = "first entrance exam Kyoto U. 2010, humanities [2]"
	return q
}

func exLass() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{1, 2, 3}, NewFmlOrs(
		NewAtom(NewPolyCoef(3, 0, 1), LT),
		NewAtom(NewPolyCoef(3, -1, 1), GE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), 1), LT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, -1), 1), LT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), 1), LT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, -1), 1), LT),
		NewAtom(NewPolyCoef(2, NewPolyCoef(0, 1, -1), 0, NewPolyCoef(1, 0, 0, -1, 0, 1), 0, NewPolyCoef(1, 0, 0, 1)), GE)))
	q.Output = NewAtom(NewPolyCoef(0, -26, 27), LE)
	q.Ref = "Jean B. Lasserre. Global Optimization with Polynomials and the Problem of Moments. 2001"
	q.DOI = "10.1137/S1052623400366802"
	return q
}

func exMakePdf() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{1}, NewFmlAnds(
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), EQ),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 1), one), LT)))
	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(0, -1, 0, 2), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(0, 0, 1), LE),
			NewAtom(NewPolyCoef(0, 1, 1), GE)))
	q.Ref = "Christopher W. Brown. Solution formula construction for truth invariant CAD's. Thesis p65 1999"
	return q
}

func exMakePdf2() *QeExample {
	// [ y^2 + x^2 - 1 <= 0 /\ 5 x + 3 < 0 ] \/ 5 y^2 + 4 x y - 4 y + 2 x^2 - 2 x < 0 \/ [ y^2 + x^2 - 1 <= 0 /\ 5 y + 2 x - 2 < 0 ]
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), 0, 1), EQ),
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, -1, 1), 2), 1), LT)))
	q.Output = NewFmlOrs(
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), LE),
			NewAtom(NewPolyCoef(0, 3, 5), LT)),
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -2, 2), NewPolyCoef(0, -4, 4), 5), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), LE),
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, -2, 2), 5), LT)))
	q.Ref = "Christopher W. Brown. Solution formula construction for truth invariant CAD's. Thesis p65 1999"
	return q
}

func exNeq1() *QeExample {
	// 非等式制約 ex([x], a*x^5+b*x^4+c*x^3+d*x^2+e*x+f != 0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{6}, NewAtom(NewPolyCoef(6, NewPolyCoef(5, 0, 1), NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE))
	q.Output = NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), NE), NewAtom(NewPolyCoef(1, 0, 1), NE), NewAtom(NewPolyCoef(2, 0, 1), NE), NewAtom(NewPolyCoef(3, 0, 1), NE), NewAtom(NewPolyCoef(4, 0, 1), NE), NewAtom(NewPolyCoef(5, 0, 1), NE))
	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exNeq2() *QeExample {
	// 非等式制約 ex([x], a*x^2+b*x + c != 0 && d*x^2+e*x+f != 0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{6}, NewFmlAnds(NewAtom(NewPolyCoef(6, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE), NewAtom(NewPolyCoef(6, NewPolyCoef(5, 0, 1), NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1)), NE)))
	q.Output = NewFmlAnds(
		NewFmlOrs(
			NewAtom(NewPolyCoef(0, 0, 1), NE),
			NewAtom(NewPolyCoef(1, 0, 1), NE),
			NewAtom(NewPolyCoef(2, 0, 1), NE)),
		NewFmlOrs(
			NewAtom(NewPolyCoef(3, 0, 1), NE),
			NewAtom(NewPolyCoef(4, 0, 1), NE),
			NewAtom(NewPolyCoef(5, 0, 1), NE)))

	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exNeq3() *QeExample {
	// まとめ (1) の改版 strict
	q := new(QeExample)
	// ex([x], x^2+a*x+b<0 && x-c!=0 && d*x^2+e*x+f!=0)
	q.Input = NewQuantifier(false, []Level{6}, NewFmlAnds(
		NewAtom(NewPolyCoef(6, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LT),
		NewAtom(NewPolyCoef(6, NewPolyCoef(2, 0, -1), 1), NE),
		NewAtom(NewPolyCoef(6, NewPolyCoef(5, 0, 1), NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1)), NE)))

	// a^2-4*b > 0 && (d != 0 || e != 0 || f != 0)
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 4), LT),
		NewFmlOrs(
			NewAtom(NewPolyCoef(3, 0, 1), NE),
			NewAtom(NewPolyCoef(4, 0, 1), NE),
			NewAtom(NewPolyCoef(5, 0, 1), NE)))

	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exNeq4() *QeExample {
	// まとめ (2) 改版 strict
	q := new(QeExample)
	// ex([x], x^2-2*x*a3+a4^2-2*a4*a2+a2^2+a3^2 < 0 && x != 0 && (a4*a1-a4*a3+a1*a2-a2*a3)*x+a4^2-a4^2*a2-a4*a1*a3+a4*a3^2 != 0);
	// discrim(x^2-2*x*a3+a4^2-2*a4*a2+a2^2+a3^2, x) = 4*(a4-a2)^2 >= 0 故，この原子論理式だけで false 確定する面白くない

	// 改版
	// ex([x], x^2-2*x*a3+a4^2-2*a4*a2+a2^2+a3^2 < 1 && x != 0 && (a4*a1-a4*a3+a1*a2-a2*a3)*x+a4^2-a4^2*a2-a4*a1*a3+a4*a3^2 != 0);
	q.Input = NewQuantifier(false, []Level{4}, NewFmlAnds(
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, -1, 0, 1), 0, 1), NewPolyCoef(1, 0, -2), 1), NewPolyCoef(2, 0, -2), 1), LT),
		NewAtom(NewPolyCoef(4, 0, 1), NE),
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, NewPolyCoef(2, 0, NewPolyCoef(0, 0, 1), -1), NewPolyCoef(1, -1, 1)), NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), NewPolyCoef(2, NewPolyCoef(0, 0, -1), 1))), NE)))

	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, -1, 0, 1), NewPolyCoef(1, 0, -2), 1), LT),
		NewFmlOrs(
			NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), NewPolyCoef(2, NewPolyCoef(0, 0, -1), 1)), NE),
			NewAtom(NewPolyCoef(3, 0, NewPolyCoef(2, 0, NewPolyCoef(0, 0, 1), -1), NewPolyCoef(1, -1, 1)), NE)))

	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exNeq5() *QeExample {
	// まとめ (1) の改版 atom
	q := new(QeExample)
	// ex([x], e*x^5+c*x^3+x^2+a*x+b<=0 && x-c!=0 && d*x^2+e*x+f!=0)
	q.Input = NewQuantifier(false, []Level{6}, NewFmlAnds(
		NewAtom(NewPolyCoef(6, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1, NewPolyCoef(2, 0, 1), 0, NewPolyCoef(4, 0, 1)), LE),
		NewAtom(NewPolyCoef(6, NewPolyCoef(2, 0, -1), 1), NE),
		NewAtom(NewPolyCoef(6, NewPolyCoef(5, 0, 1), NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1)), NE)))

	// (e != 0 || c != 0 || a^2-4*b > 0) && (d != 0 || e != 0 || f != 0) ||
	// (e == 0 && c == 0 && a^2-4*b == 0 && ex([x], 2*x+a== 0 && x!=c && d*x^2+e*x+f != 0)
	q.Output = NewFmlOrs(
		NewFmlAnds(
			NewFmlOrs(
				NewAtom(NewPolyCoef(4, 0, 1), NE),
				NewAtom(NewPolyCoef(2, 0, 1), NE),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 4), LT)),
			NewFmlOrs(
				NewAtom(NewPolyCoef(3, 0, 1), NE),
				NewAtom(NewPolyCoef(4, 0, 1), NE),
				NewAtom(NewPolyCoef(5, 0, 1), NE))),
		NewFmlAnds(
			NewAtom(NewPolyCoef(4, 0, 1), EQ),
			NewAtom(NewPolyCoef(2, 0, 1), EQ),
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 4), EQ),
			NewAtom(NewPolyCoef(0, 0, 1), NE),
			NewAtom(NewPolyCoef(5, NewPolyCoef(3, 0, NewPolyCoef(0, 0, 0, 1)), 4), NE)))

	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exNeq6() *QeExample {
	// まとめ (3) strictに近い
	q := new(QeExample)
	// ex([x], x<0 && (a3-3)*x+2>0 && (a3^2-8*a3+24)*x+4*a3-20>=0 && a1*x-a2!=0 && (a1*a3+a1^2-4*a1-2)*x^2+(-a2*a3+(-2*a1+4)*a2+2*a1+1)*x+a2^2-2*a2+1!=0)
	q.Input = NewQuantifier(false, []Level{4},
		NewFmlAnds(
			NewAtom(NewPolyCoef(4, 0, 1), LT),
			NewAtom(NewPolyCoef(4, 2, NewPolyCoef(2, -3, 1)), GT),
			NewAtom(NewPolyCoef(4, NewPolyCoef(2, -20, 4), NewPolyCoef(2, 24, -8, 1)), GE),
			NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, -1), NewPolyCoef(0, 0, 1)), NE),
			NewAtom(NewPolyCoef(4, NewPolyCoef(1, 1, -2, 1), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 1, 2), NewPolyCoef(0, 4, -2)), NewPolyCoef(1, 0, -1)), NewPolyCoef(2, NewPolyCoef(0, -2, -4, 1), NewPolyCoef(0, 0, 1))), NE)))

	// (!)  = a1*x-a2!=0 && (a1*a3+a1^2-4*a1-2)*x^2+(-a2*a3+(-2*a1+4)*a2+2*a1+1)*x+a2^2-2*a2+1!=0
	// (!') = ex([x], (!))
	//      = (a1 != 0 || a2 != 0) && (a1*a3+a1^2-4*a1-2 != 0 || -a2*a3+(-2*a1+4)*a2+2*a1+1 != 0 || a2^2-2*a2+1!=0)
	//   ex([x], x<0 && (a3-3)*x+2>0 && (a3^2-8*a3+24)*x+4*a3-20>0) && ex([x], (!)) ||
	//   ex([x], x<0 && (a3-3)*x+2>0 && (a3^2-8*a3+24)*x+4*a3-20=0 && (!))
	// <==>
	//   a3 > 5 && (!') ||
	//   Q<0 && (a3-3)*Q+2*P > 0 && a1*Q-a2*P != 0 && (a1*a3+a1^2-4*a1-2)*Q^2+(-a2*a3+(-2*a1+4)*a2+2*a1+1)*P*Q+(a2^2-2*a2+1)*Q!=0
	//   where P = a3^2-8*a3+24, Q = -(4*a3-20)
	//
	// note that: a3^2-8*a3+24 > 0 by the discriminant is negative
	q.Output = NewFmlOrs(
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, -20, 4), GT),
			NewAtom(NewPolyCoef(2, 12, -16, 2), LT),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -20), 24), NewPolyCoef(1, NewPolyCoef(0, 0, 4), -8), NewPolyCoef(1, 0, 1)), NE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 80, -640, 400), NewPolyCoef(0, 1120, -960), 400), NewPolyCoef(1, NewPolyCoef(0, -96, 528, -160), NewPolyCoef(0, -1184, 512), -160), NewPolyCoef(1, NewPolyCoef(0, 36, -120, 16), NewPolyCoef(0, 432, -104), 16), NewPolyCoef(1, NewPolyCoef(0, -4, 8), NewPolyCoef(0, -68, 8)), NewPolyCoef(1, 0, 4)), NE)),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, -5, 1), GT),
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 0, 1), NE),
				NewAtom(NewPolyCoef(1, 0, 1), NE)),
			NewFmlOrs(
				NewAtom(NewPolyCoef(2, NewPolyCoef(0, -2, -4, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, -1, -2), NewPolyCoef(0, -4, 2)), NewPolyCoef(1, 0, 1)), NE),
				NewAtom(NewPolyCoef(1, -1, 1), NE))))

	q.Ref = "Iwane. Quantifier elimination for inequational constraints. 2015"
	return q
}

func exPL01() *QeExample {
	q := new(QeExample)
	q.Output = NewAtom(NewPolyCoef(0, 0, 1), LE)
	q.Input = NewQuantifier(true, []Level{1, 2}, NewFmlOrs(
		NewAtom(NewPolyCoef(1, 1, 1), LT),
		NewAtom(NewPolyCoef(1, -1, 1), GT),
		NewAtom(NewPolyCoef(2, 1, 1), LT),
		NewAtom(NewPolyCoef(2, -1, 1), GT),
		NewAtom(NewPolyCoef(2, NewPolyCoef(0, 1, -1), 0, NewPolyCoef(1, 0, 0, -3, 0, 1), 0, NewPolyCoef(1, 0, 0, 1)), GE)))
	q.Ref = "P. Parrilo and S. Lall. Semidefinite Programming Relaxation and Algebraic Optimization in Control."
	q.DOI = ""
	return q
}

func exPortfolio() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{1, 2, 3}, NewFmlAnds(
		NewAtom(NewPolyCoef(1, 0, 1), GE),
		NewAtom(NewPolyCoef(2, 0, 1), GE),
		NewAtom(NewPolyCoef(3, 0, 1), GE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, -10000, 1), 1), 1), LE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, -100000, 5), -4), 15), GE),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, 8), NewPolyCoef(1, 0, -40), 3), NewPolyCoef(2, NewPolyCoef(1, 0, 10), -30), 45), EQ)))
	q.Ref = "J. Schattman. Portfolio optimization under nonconvex transaction costs with the global optimization toolbox. https://www.maplesoft.com/Applications/Detail.aspx?id=1401"
	return q
}

func exQuad() *QeExample {
	// ex([x], a*x^2+b*x+c == 0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), EQ))
	q.Output = NewFmlOrs(NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, -1), NewPolyCoef(0, 0, 4)), LE),
		NewAtom(NewPolyCoef(0, 0, 1), NE)), NewFmlAnds(
		NewAtom(NewPolyCoef(0, 0, 1), EQ),
		NewAtom(NewPolyCoef(1, 0, 1), NE)),
		NewAtom(NewPolyCoef(2, 0, 1), EQ))

	return q
}

func exQuart() *QeExample {
	// all([x], x^4+p*x^2+q*x+r>=0)
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{3}, NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 0, 1), GE))

	// 256 r^3 - 128 p^2 r^2 + 144 p q^2 r + 16 p^4 r - 27 q^4 - 4 p^3 q^2 >= 0 /\ [ 27 q^2 + 8 p^3 > 0 \/ [ 48 r^2 - 16 p^2 r + 9 p q^2 + p^4 >= 0 /\ 6 r - p^2 >= 0 ] ]
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, 0, -4), 0, -27), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 16), 0, NewPolyCoef(0, 0, 144)), NewPolyCoef(0, 0, 0, -128), 256), GE),
		NewFmlOrs(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 8), 0, 27), GT),
			NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 0, 9)), NewPolyCoef(0, 0, 0, -16), 48), GE),
				NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 0, -1), 6), GE))))
	return q
}

func exSato1() *QeExample {
	// ex([y], y^6-2*y^4-(a+2)*y^3+y^2+2*y+1==0);
	// <==>
	// 729*a^3+8856*a^2+43632*a+33856<=0 || a>=0
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{1}, NewAtom(NewPolyCoef(1, 1, 2, 1, NewPolyCoef(0, -2, -1), -2, 0, 1), EQ))
	q.Output = NewFmlOrs(NewAtom(NewPolyCoef(0, 33856, 43632, 8856, 729), LE), NewAtom(NewPolyCoef(0, 0, 1), GE))
	return q
}

func exSDC2() *QeExample {
	// ex([x], x>=0 && x^2+b*x+c <= 0)
	// <==>
	// all([x], x >= 0 => x^2+b*x+c > 0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2}, NewFmlAnds(NewAtom(NewPolyCoef(2, 0, 1), GE), NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LE)))
	// 4*c-b^2<=0 && (b<0 || c<=0)
	q.Output = NewFmlAnds(NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 4), LE), NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), LT), NewAtom(NewPolyCoef(1, 0, 1), LE)))
	q.Ref = "An Effective Implementation of a Special Quantifier Elimination for a Sign Definite Condition by Logical Formula Simplification"
	q.DOI = "https://doi.org/10.1007/978-3-319-02297-0_17"
	return q
}

func exRoot2() *QeExample {
	// ex([x], x^2+b*x+c <= 0)
	// <==>
	// all([x], x^2+b*x+c > 0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2}, NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LE))
	// 4*c-b^2<=0 && (b<0 || c<=0)
	q.Output = NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), 4), LE)
	q.Ref = "Gonzalez-Vega, Laureano.  A combinatorial algorithm solving some quantifier elimination problems.  1998"
	return q
}

func exSDC3() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewFmlAnds(NewAtom(NewPolyCoef(3, 0, 1), GE), NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LE)))

	d := NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, 1), -4), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -4), NewPolyCoef(0, 0, 18)), -27)

	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(2, 0, 1), LE),
		NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), LT), NewAtom(d, GE)),
		NewFmlAnds(NewAtom(NewPolyCoef(1, 0, 1), LT), NewAtom(d, GE)))
	q.Ref = "An Effective Implementation of a Special Quantifier Elimination for a Sign Definite Condition by Logical Formula Simplification"
	q.DOI = "https://doi.org/10.1007/978-3-319-02297-0_17"

	return q
}

func exRoot3() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LE))
	q.Output = trueObj
	q.Ref = "Gonzalez-Vega, Laureano.  A combinatorial algorithm solving some quantifier elimination problems.  1998"
	return q
}

func exSDC4() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{4}, NewFmlAnds(NewAtom(NewPolyCoef(4, NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(3, 0, 1), 1), LE), NewAtom(NewPolyCoef(4, 0, 1), GE)))
	q.Output = nil
	q.Ref = "An Effective Implementation of a Special Quantifier Elimination for a Sign Definite Condition by Logical Formula Simplification"
	q.DOI = "https://doi.org/10.1007/978-3-319-02297-0_17"
	return q
}

func exRoot4() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{4}, NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1), 1), LE))

	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, 1), -4), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -4), NewPolyCoef(0, 0, 18)), -27), NewPolyCoef(2, NewPolyCoef(1, 0, 0, 0, NewPolyCoef(0, 0, 0, -4), 16), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 18), NewPolyCoef(0, 0, -80)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -6), 144)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, -27), NewPolyCoef(0, 0, 0, 144), -128), NewPolyCoef(0, 0, -192)), 256), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, 1), -4), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -4), NewPolyCoef(0, 0, 18)), -27), NewPolyCoef(2, NewPolyCoef(1, 0, 0, 0, NewPolyCoef(0, 0, 0, -4), 16), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 18), NewPolyCoef(0, 0, -80)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -6), 144)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, -27), NewPolyCoef(0, 0, 0, 144), -128), NewPolyCoef(0, 0, -192)), 256), EQ),
			NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, 0, NewPolyCoef(0, 0, 0, -2), 8), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 9), NewPolyCoef(0, 0, -40)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -3), 72)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, -27), NewPolyCoef(0, 0, 0, 144), -128), NewPolyCoef(0, 0, -192)), 384), GT)),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, -9), 32), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 27), NewPolyCoef(0, 0, -108)), 108), LE),
			NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, 0, NewPolyCoef(0, 0, 0, -2), 8), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 9), NewPolyCoef(0, 0, -40)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -3), 72)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, -27), NewPolyCoef(0, 0, 0, 144), -128), NewPolyCoef(0, 0, -192)), 384), LE)),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, -9), 32), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 27), NewPolyCoef(0, 0, -108)), 108), LE),
			NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, -27), NewPolyCoef(0, 0, 0, 144), -128), NewPolyCoef(0, 0, -192)), 768), LT)))

	q.Ref = "Gonzalez-Vega, Laureano.  A combinatorial algorithm solving some quantifier elimination problems.  1998"
	return q
}

func exRoot4_2() *QeExample {
	// ex([x], a*x^4+b*x^2+c*x+1<=0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, 1, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), LE))
	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -256), 0, NewPolyCoef(0, 0, 128), 0, -16), 0, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -144), 0, 4), 0, NewPolyCoef(0, 0, 27)), GT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, 0, 1), NE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -256), 0, NewPolyCoef(0, 0, 128), 0, -16), 0, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -144), 0, 4), 0, NewPolyCoef(0, 0, 27)), EQ),
		),
		NewAtom(NewPolyCoef(0, 0, 1), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -4), 0, 1), GE),
			NewAtom(NewPolyCoef(1, 0, 1), LT)))
	return q
}

func exRoot4_3() *QeExample {
	// ex([x], a*x^4+b*x^2+c*x+1==0)
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, 1, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), EQ))
	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -256), 0, NewPolyCoef(0, 0, 128), 0, -16), 0, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -144), 0, 4), 0, NewPolyCoef(0, 0, 27)), GT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, 0, 1), NE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -256), 0, NewPolyCoef(0, 0, 128), 0, -16), 0, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -144), 0, 4), 0, NewPolyCoef(0, 0, 27)), EQ),
		),
		NewAtom(NewPolyCoef(0, 0, 1), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -4), 0, 1), GE),
			NewAtom(NewPolyCoef(1, 0, 1), LT)))

	return q
}

func exWlog1() *QeExample {
	/* 正三角形のひとつの角の余弦
	 all([x,y,z,w,a,b], impl((x != z || y != w) &&
		(x - z)^2 + (y - w)^2 == (z - a)^2 + (w - b)^2 &&
		(z - a)^2 + (w - b)^2 == (a - x)^2 + (b - y)^2,
		(z - a)^2 + (w - b)^2 == (a - x)^2 + (b - y)^2 + (x - z)^2 + (y - w)^2 - 2 * c
	*/
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{1, 2, 3, 4, 5, 6}, NewFmlOrs(
		NewFmlAnds(
			NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, -1), 1), EQ),
			NewAtom(NewPolyCoef(4, NewPolyCoef(2, 0, -1), 1), EQ)),
		NewAtom(NewPolyCoef(6, NewPolyCoef(5, NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, -1), 0, -1), NewPolyCoef(1, 0, 2)), NewPolyCoef(2, 0, 2)), NewPolyCoef(3, 0, -2), 1), NewPolyCoef(4, 0, -2), 1), NE),
		NewAtom(NewPolyCoef(6, NewPolyCoef(5, NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, 1), 0, 1), 0, -1), 0, -1), NewPolyCoef(3, NewPolyCoef(1, 0, -2), 2)), NewPolyCoef(4, NewPolyCoef(2, 0, -2), 2)), NE),
		NewAtom(NewPolyCoef(6, NewPolyCoef(5, NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, -2, 2)), 0, NewPolyCoef(0, -2, 2)), NewPolyCoef(1, 0, 2)), NewPolyCoef(2, 0, 2)), NewPolyCoef(3, NewPolyCoef(1, 0, NewPolyCoef(0, 2, -4)), -2), NewPolyCoef(0, 0, 2)), NewPolyCoef(4, NewPolyCoef(2, 0, NewPolyCoef(0, 2, -4)), -2), NewPolyCoef(0, 0, 2)), EQ)))
	q.Output = NewAtom(NewPolyCoef(0, -1, 2), EQ)
	q.Ref = "H. Iwane, H. Anai. Formula Simplification for Real Quantifier Elimination Using Geometric Invariance"
	q.DOI = "10.1145/3087604.3087627"
	return q
}

func exWO1() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{3}, NewAtom(NewPolyCoef(3, 1, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), 0, 0, 0, NewPolyCoef(0, 0, 1)), GT))
	q.Output = q.Input
	q.Ref = "original: NOT well-oriented"
	return q
}

func exWO2() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{3}, NewAtom(
		NewPolyCoef(3,
			NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)),
			NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1)),
			1), GE))
	q.Output = NewFmlOrs(
		NewFmlAnds(NewAtom(NewPolyCoef(1, 0, 1), GE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 4), 0, NewPolyCoef(1, 0, 1)), LE)),
		NewFmlAnds(NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 4), 0, NewPolyCoef(1, 0, 1)), GE),
			NewAtom(NewPolyCoef(1, 0, 1), LE)))
	q.Ref = "original: NOT well-oriented, but no delineating polynomial is needed!"
	return q
}

func exWO3() *QeExample {
	// (c,d,b,x)
	// 3
	// (E x) [ x >= 0 /\ x^3 + b x^2 + c x + d < 0 ].
	// Error! Delineating polynomial should be added over cell(2,2)!
	// d-cell (2,2) -> (x=0, y=0)
	// Degrees after substitution  : (-1)
	// x=y=0
	// P_3,1  = fac(J_3,1) = fac(dis(A_4,1))
	//        = 4 c^3 - b^2 c^2 - 18 b d c + 27 d^2 + 4 b^3 d
	// x <= 0 \/ [ 4 x z^3 - y^2 z^2 - 18 x y z + 4 y^3 + 27 x^2 <= 0 /\ 6 x z^2 - y^2 z - 9 x y > 0 /\ 12 x z - y^2 < 0 ] \/ [ y^3 - 27 x^2 < 0 /\ 4 x z^3 - y^2 z^2 - 18 x y z + 4 y^3 + 27 x^2 <= 0 ]
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{3}, NewFmlAnds(NewAtom(NewPolyCoef(3, 0, 1), GE), NewAtom(NewPolyCoef(3, NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1), 1), LT)))
	q.Output = NewFmlOrs(
		NewAtom(NewPolyCoef(0, 0, 1), LT),
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, 216), 0, 0, 1), EQ),
			NewAtom(NewPolyCoef(2, 0, 1), LT)),
		NewFmlAnds(
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 0, -27), 0, 0, 1), LT),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 27), 0, 0, 4), NewPolyCoef(1, 0, NewPolyCoef(0, 0, -18)), NewPolyCoef(1, 0, 0, -1), NewPolyCoef(0, 0, 4)), LT)),
		NewFmlAnds(
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 27), 0, 0, 4), NewPolyCoef(1, 0, NewPolyCoef(0, 0, -18)), NewPolyCoef(1, 0, 0, -1), NewPolyCoef(0, 0, 4)), LT),
			NewAtom(NewPolyCoef(2, 0, 1), LT)))
	q.Ref = "original(SDC): NOT well-oriented"
	return q
}

func exWO4() *QeExample {
	// (x,y,z,w)
	// 3
	// (E w) [ w^2 < x /\ z w + y <= 0 ].
	q := new(QeExample)
	q.Ref = "original: NOT well-oriented, but no delineating polynomial is needed!"
	q.Input = NewQuantifier(false, []Level{3}, NewFmlAnds(
		NewAtom(NewPolyCoef(3, NewPolyCoef(0, 0, -1), 0, 1), LT),
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1)), LE)))
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(0, 0, 1), GT),
		NewFmlOrs(
			NewAtom(NewPolyCoef(1, 0, 1), LE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, -1), 0, NewPolyCoef(0, 0, 1)), GT)))
	return q
}

func exXAxisEllipse() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{3, 4}, NewFmlAnds(
		NewAtom(NewPolyCoef(0, 0, 1), GT),
		NewAtom(NewPolyCoef(1, 0, 1), GT), NewFmlOrs(
			NewAtom(NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, 0, -1)), 0, NewPolyCoef(1, 0, 0, 1)), NewPolyCoef(2, 0, NewPolyCoef(1, 0, 0, -2)), NewPolyCoef(1, 0, 0, 1)), 0, NewPolyCoef(0, 0, 0, 1)), NE),
			NewAtom(NewPolyCoef(4, NewPolyCoef(3, -1, 0, 1), 0, 1), LE))))
	q.Output = q.Input
	q.Ref = "The x-axix ellipse problem: W. Kahan. Problem no. 9: An ellipse problem."
	return q
}

func exWhitneyUmbrella() *QeExample {
	q := new(QeExample)
	// ex([u,v], u*v-x==0 && v-y==0 && u^2-z==0)
	q.Input = NewQuantifier(false, []Level{3, 4}, NewFmlAnds(
		NewAtom(NewPolyCoef(4, NewPolyCoef(0, 0, -1), NewPolyCoef(3, 0, 1)), EQ),
		NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, -1), 1), EQ),
		NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, -1), 0, 1), EQ)))
	q.Output = NewFmlAnds(
		NewAtom(NewPolyCoef(2, 0, 1), GE),
		NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 0, -1), NewPolyCoef(1, 0, 0, 1)), EQ))
	return q
}

func exHong93() *QeExample {
	q := new(QeExample)
	// vars(u,v,w,x);
	// ex([x], u*x^2+v*x+1==0 && v*x^3+w*x+u==0 && w*x^2+v*x+u <= 0);
	q.Input = NewQuantifier(false, []Level{3}, NewFmlAnds(
		NewAtom(NewPolyCoef(3, 1, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), EQ),
		NewAtom(NewPolyCoef(3, NewPolyCoef(0, 0, 1), NewPolyCoef(2, 0, 1), 0, NewPolyCoef(1, 0, 1)), EQ),
		NewAtom(NewPolyCoef(3, NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1)), LE)))

	// RB=u^5-w*v*u^3+(3*v^2+w^2)*u^2+(-v^4-2*w*v)*u+w*v^3+v^2;
	// TB=2*u^4-w*v*u^2+3*v^2*u-v^4;
	// SB=w*u^2-v*u+v^3;
	// RC=u^4+(-v^2-2*w)*u^2+(w+1)*v^2*u-w*v^2+w^2;
	// TC=2*u^3+(-v^2-2*w)*u+w*v^2;
	// SC=v*u-w*v;
	q.Output = NewFmlOrs(
		// u=0 && v != 0 && F'
		NewFmlAnds(
			NewAtom(NewPolyCoef(0, 0, 1), EQ),
			NewAtom(NewPolyCoef(1, 0, 1), NE),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 0, 1)), EQ),
			NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 0, NewPolyCoef(0, -1, 1)), 1), LE)),
		NewFmlAnds(
			// u != 0 && v^2-4*u >= 0
			NewAtom(NewPolyCoef(0, 0, 1), NE),
			NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -4), 0, 1), GE),
			NewFmlOrs(
				// F1
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 1, 0, 3), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, -2, 0, -1), 0, 1), NewPolyCoef(0, 0, 0, 1)), EQ),
					NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, 0, 3, 0, -2), 0, NewPolyCoef(0, 0, -4), 0, 1), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 0, 0, -2), 0, NewPolyCoef(0, 0, 0, 0, -4), 0, NewPolyCoef(0, 0, 0, 2)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 0, 1))), GE),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 0, 1, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2), 0, NewPolyCoef(0, -1, 1)), 1), LE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, -2), 0, 1)), LE)),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 0, 1, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2), 0, NewPolyCoef(0, -1, 1)), 1), GE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), GE)),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, -2), 0, 1)), LE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), GE))),
				// F2
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 1, 0, 3), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, -2, 0, -1), 0, 1), NewPolyCoef(0, 0, 0, 1)), EQ),
					NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, 0, 3, 0, -2), 0, NewPolyCoef(0, 0, -4), 0, 1), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 0, 0, -2), 0, NewPolyCoef(0, 0, 0, 0, -4), 0, NewPolyCoef(0, 0, 0, 2)), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 0, 0, 0, 1))), LE),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 0, 1, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2), 0, NewPolyCoef(0, -1, 1)), 1), LE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, -2), 0, 1)), LE)),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), 0, NewPolyCoef(0, 0, 1, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, 0, -2), 0, NewPolyCoef(0, -1, 1)), 1), GE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), LE)),
					NewFmlOrs(
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 2), 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, NewPolyCoef(0, 0, -2), 0, 1)), LE),
						NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, -1)), NewPolyCoef(1, 0, 1)), LE))))))

	q.Ref = "Hoon Hong. Quantifier elimination for formulas constrained by quadratic equations via slope resultants"
	q.DOI = "10.1145/164081.164140"
	return q
}

func exImo13_1_5() *QeExample {
	// all([a1,a2,a3,a4,a5], (a1-a2)*(a1-a3)*(a1-a4)*(a1-a5)+(a2-a1)*(a2-a3)*(a2-a4)*(a2-a5)+(a3-a1)*(a3-a2)*(a3-a4)*(a3-a5)+(a4-a1)*(a4-a2)*(a4-a3)*(a4-a5)+(a5-a1)*(a5-a2)*(a5-a3)*(a5-a4) >= 0);
	q := new(QeExample)
	q.Input = NewQuantifier(true, []Level{0, 1, 2, 3, 4}, NewAtom(
		NewPolyCoef(4,
			NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, 0, 1), NewPolyCoef(0, 0, 0, 0, -1), 0, NewPolyCoef(0, 0, -1), 1), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -1), NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, 1), -1), NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1)), NewPolyCoef(1, NewPolyCoef(0, 0, -1), -1), 1), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -1), NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, 1), -1), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, -3), 1), NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), -1), NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1)), NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), -1), -1), 1),
			NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 0, -1), NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, 1), -1), NewPolyCoef(1, NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, -3), 1), NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), -1), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, 1), NewPolyCoef(0, 0, -3), 1), NewPolyCoef(1, NewPolyCoef(0, 0, -3), -3), 1), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), 1), -1),
			NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1)), NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1)), NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), 1)),
			NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -1), -1), -1), -1),
			1), GE))
	q.Output = trueObj
	q.Ref = "H. Iwane, H. Anai. Formula Simplification for Real Quantifier Elimination Using Geometric Invariance"
	q.DOI = "10.1145/3087604.3087627"
	return q
}

func exMooea() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{2, 3, 4}, NewFmlAnds(
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(0, 0, -1), 0, 1), 0, 1), 1), EQ),
		NewAtom(NewPolyCoef(4, NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, 1, -1), -2, 1), 0, 1), 1), EQ),
		NewAtom(NewPolyCoef(2, 2, 1), GE), NewAtom(NewPolyCoef(2, -2, 1), LE),
		NewAtom(NewPolyCoef(3, 2, 1), GE), NewAtom(NewPolyCoef(3, -2, 1), LE), NewAtom(NewPolyCoef(4, 1, 10), GE), NewAtom(NewPolyCoef(4, -1, 10), LE)))
	return q
}

func exWilson() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{0, 1}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, -5, 4, -1), 2, -1), 1), EQ),
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, NewPolyCoef(0, -36, 0, -1), 12, -1), 1), EQ),
		NewAtom(NewPolyCoef(0, -2, 5), GE),
		NewAtom(NewPolyCoef(0, -8, 5), LE),
		NewAtom(NewPolyCoef(1, -2, 1), GE),
		NewAtom(NewPolyCoef(1, -5, 1), LE)))
	return q
}

func exLampinen() *QeExample {
	q := new(QeExample)
	q.Input = NewQuantifier(false, []Level{0, 1}, NewFmlAnds(
		NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -1), -1), 1), EQ),
		NewAtom(NewPolyCoef(3, NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, -1), 1), EQ),
		NewAtom(NewPolyCoef(0, 10, 1), GE),
		NewAtom(NewPolyCoef(0, -10, 1), LE),
		NewAtom(NewPolyCoef(1, 10, 1), GE),
		NewAtom(NewPolyCoef(1, -10, 1), LE)))
	return q
}
