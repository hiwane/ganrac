package main

// espresso コマンドの出力 pla ファイルを解析する

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hiwane/ganrac"
	"os"
)

// sgn_table.delta() のコピー
func shDelta(n int) int {
	return -1 + ((n - 1) & 2)
}

// fof OP.neg() のコピー
func opNeg(op ganrac.OP) ganrac.OP {
	if op == ganrac.EQ || op == ganrac.NE || op == ganrac.OP_TRUE || op == ganrac.OP_FALSE {
		return op
	}
	return op ^ (ganrac.LT | ganrac.GT)
}

type Esp struct {
	verbose bool
	src     int
	prefix  string
	cost    int
	onset   [][]ganrac.OP
	offset  [][]ganrac.OP
}

func (esp *Esp) dprint(format string, a ...any) {
	if esp.verbose {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func (esp *Esp) Cost(opp [][]ganrac.OP) int {
	if len(opp) <= 0 || esp.cost <= 0 {
		return -1
	}
	n := len(opp[0])
	w := make([]int, n)
	base := 2
	if esp.cost == 1 {
		// SDC.  Sturm-Habitch列は，集結式が一番次数が高くなる
		for i := 0; i < n/2; i++ {
			w[i] = i + base + 1
			w[i+n/2] = i + base
		}
	} else if esp.cost == 2 {
		for i := 0; i < n; i++ {
			w[i] = i + base
		}
	} else {
		for i := 0; i < n; i++ {
			w[i] = base
		}
	}

	for _, ww := range w {
		if ww <= 1 {
			panic(fmt.Sprintf("all element of w must be positive: %v", w))
		}
	}

	ret := 0
	for _, ops := range opp {
		for i, o := range ops {
			if o == ganrac.OP_TRUE {
				continue
			} else if o == ganrac.EQ {
				ret++
			} else {
				ret += w[i]
			}
		}
	}
	return ret
}

func (esp *Esp) PrintSrc(opp [][]ganrac.OP, tab string) {
	fmt.Printf("%s{	// ", tab)
	n := len(opp[0])
	if esp.src == 1 { // sdc
		fmt.Printf("SDC mode\n")
		n = n / 2
	} else if esp.src == 2 { // sdc
		fmt.Printf("ATOM mode\n")
		n = n / 2
	} else {
		fmt.Printf("normal mode\n")
	}

	for _, ops := range opp {
		fmt.Printf("%s\t", tab)
		sep := "{"
		for i, o := range ops {
			var j = -1
			if esp.src == 1 { // sdc
				if i < n {
					j = i + 1
				} else if i != n {
					j = i - n
				}
			} else if esp.src == 2 { // atom
				j = i + 1
			}
			if j >= 0 && shDelta(j) < 0 {
				// SH列は，部分集結式列の符号を変えたものだが，
				// 呼び出し元は部分集結式列を計算し，そのまま使えるようにするため，
				// ここで符号をいじる
				o = opNeg(o)
			}

			fmt.Printf("%s%s%S", sep, esp.prefix, o)
			sep = ", "
		}
		fmt.Printf("},\n")
	}
	fmt.Printf("%s}\n", tab)
}

// PLA 形式のファイルを読み込み.
// true 行のみを読み込む
// returns (ONset, OFFset)
func parse(fp *os.File) ([][]ganrac.OP, [][]ganrac.OP) {
	scanner := bufio.NewScanner(fp)
	onset := make([][]ganrac.OP, 0)
	offset := make([][]ganrac.OP, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '.' || line[0] == 's' || (line[len(line)-1] != '1' && line[len(line)-1] != '0') {
			continue
		}
		rr := make([]ganrac.OP, 0, 20)
		for i := 0; i < len(line)-2; i += 2 {
			var op ganrac.OP
			if line[i] == '-' && line[i+1] == '-' {
				op = ganrac.OP_TRUE
			} else if line[i] == '1' && (line[i+1] == '-' || line[i+1] == '0') {
				op = ganrac.GT
			} else if line[i+1] == '1' && (line[i] == '-' || line[i] == '0') {
				op = ganrac.LT
			} else if line[i] == '0' && line[i+1] == '-' { // 00 01
				op = ganrac.LE
			} else if line[i] == '-' && line[i+1] == '0' { // 00 10
				op = ganrac.GE
			} else if line[i] == '0' && line[i+1] == '0' { // 00
				op = ganrac.EQ
			} else {
				panic("i dont know: `" + line[i:i+2] + "` :" + line)
			}
			rr = append(rr, op)
		}
		if line[len(line)-1] == '1' {
			onset = append(onset, rr)
		} else {
			offset = append(offset, rr)
		}
	}
	return onset, offset
}

func (esp *Esp) capture(opp, input [][]ganrac.OP) bool {
	for _, in := range input {
		capt := false
		for _, ops := range opp {
			capt = true
			for i := 0; i < len(ops); i++ {
				if (ops[i] & in[i]) != in[i] {
					capt = false
					break
				}
			}
			if capt {
				// fmt.Fprintf(os.Stderr, "%v is captured by %v\n", in, ops)
				break
			}
		}
		if !capt {
			// fmt.Fprintf(os.Stderr, "== capture false\n")
			return false
		}
	}
	// fmt.Fprintf(os.Stderr, "== capture true\n")
	return true
}

func (esp *Esp) simplify(opp [][]ganrac.OP) [][]ganrac.OP {
	if !esp.capture(opp, esp.onset) {
		panic("nandeyo")
	}
	neq := 0
	fmt.Fprintf(os.Stderr, "simplify()\n")
	for i := 0; i < len(opp); i++ {
		for j := 0; j < len(opp[i]); j++ {
			if opp[i][j] != ganrac.OP_TRUE && opp[i][j]&ganrac.EQ != 0 && opp[i][j] != ganrac.EQ {
				bk := opp[i][j]
				opp[i][j] = ganrac.EQ
				if bk == opp[i][j] || (bk&opp[i][j]) != opp[i][j] || (bk|opp[i][j]) != bk {
					panic("nazo1")
				}
				// fmt.Fprintf(os.Stderr, "opp[%d][%d] = EQ\n", i, j)
				if !esp.capture(opp, esp.onset) {
					opp[i][j] = bk
				} else {
					neq++
					//	esp.dprint("simplified[%d][%d] %x => %x\n", i, j, bk, EQ)
				}
			}
		}
	}
	nne := 0
	for i := 0; i < len(opp); i++ {
		for j := 0; j < len(opp[i]); j++ {

			if opp[i][j] != ganrac.OP_TRUE && opp[i][j]&ganrac.EQ != 0 && opp[i][j] != ganrac.EQ {
				bk := opp[i][j]
				opp[i][j] &= ^ganrac.EQ
				// fmt.Fprintf(os.Stderr, "opp[%d][%d] = %v\n", i, j, opp[i][j])
				if bk == opp[i][j] || (bk&opp[i][j]) != opp[i][j] || (bk|opp[i][j]) != bk {
					panic("nazo2")
				}
				if !esp.capture(opp, esp.onset) {
					opp[i][j] = bk
				} else {
					nne++
					//	esp.dprint("simplified[%d][%d] %v => %v\n", i, j, bk, opp[i][j])
				}
			}
		}
	}

	ntr := 0
	nfa := 0
	if len(esp.offset) > 0 && false {
		// espresso がまともなら, ここが有効になるわけがない
		for i := 0; i < len(opp); i++ {
			for j := 0; j < len(opp[i]); j++ {
				if opp[i][j] != ganrac.OP_TRUE {
					bk := opp[i][j]
					opp[i][j] = ganrac.OP_TRUE
					fmt.Fprintf(os.Stderr, "opp[%d][%d] = %d => ltop\n", i, j, bk)
					for _, off := range esp.offset {
						if esp.capture(opp, [][]ganrac.OP{off}) {
							opp[i][j] = bk
							break
							//	esp.dprint("simplified[%d][%d] %v => %v\n", i, j, bk, opp[i][j])
						}
					}
					if opp[i][j] != bk {
						fmt.Fprintf(os.Stderr, "updated! [%d,%d] ltop\n", i, j)
						ntr++
					}
				}
				if opp[i][j] == ganrac.LT || opp[i][j] == ganrac.GT {
					bk := opp[i][j]
					opp[i][j] = ganrac.NE
					fmt.Fprintf(os.Stderr, "opp[%d][%d] = NE\n", i, j)
					for _, off := range esp.offset {
						if esp.capture(opp, [][]ganrac.OP{off}) {
							opp[i][j] = bk
							break
							//	esp.dprint("simplified[%d][%d] %v => %v\n", i, j, bk, opp[i][j])
						}
					}
					if opp[i][j] != bk {
						fmt.Fprintf(os.Stderr, "updated! [%d,%d] NE\n", i, j)
						nfa++
						//	esp.dprint("simplified[%d][%d] %v => %v\n", i, j, bk, opp[i][j])
					}
				}
			}
		}
	}
	esp.dprint("simplified toEQ=%d, delEQ=%d, ntr=%d,nfa=%d (%d,%d)\n", neq, nne, ntr, nfa, len(esp.onset), len(esp.offset))
	return opp
}

func (esp *Esp) gooooo(fname string, fp *os.File) {

	if fp == nil {
		var err error
		fp, err = os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "os.Open(%s) error: %v\n", fname, err)
		}
		defer fp.Close()
		esp.dprint("open file: %s\n", fname)
	}

	opp, _ := parse(fp)
	if esp.onset != nil && len(esp.onset) > 0 {
		old := len(opp)
		opp = esp.simplify(opp)
		esp.dprint("simplify: %d -> %d\n", old, len(opp))
	}

	if esp.cost > 0 {
		c := esp.Cost(opp)
		fmt.Printf("%s: cost=%d, %d\n", fname, len(opp), c)
	} else {
		fmt.Fprintf(os.Stderr, "%s: arg.cost=%d\n", fname, esp.cost)
	}

	if esp.src != 0 {
		esp.PrintSrc(opp, "\t")
	}
}

func main() {
	var esp Esp
	flag.IntVar(&esp.cost, "cost", 0, "print cost; 0: none, 1: sdc, 2: atom")
	flag.StringVar(&esp.prefix, "prefix", "", "prefix of OP")
	flag.BoolVar(&esp.verbose, "v", false, "verbose")
	flag.IntVar(&esp.src, "src", 0, "1: SDC, 2, ATOM")
	var in = flag.String("in", "", "input file for simplify")

	flag.Parse()

	if *in != "" {
		fin, err := os.Open(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "os.Open(%s) error: %v\n", *in, err)
		}

		esp.onset, esp.offset = parse(fin)
		fin.Close()
	}

	if flag.NArg() == 0 {
		esp.dprint("read stdin\n")
		esp.gooooo("", os.Stdin)
	} else {
		for _, f := range flag.Args() {
			esp.gooooo(f, nil)
		}
	}

}
