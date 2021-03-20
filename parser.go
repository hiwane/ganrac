// Code generated by goyacc -o parser.go parser.y. DO NOT EDIT.

//line parser.y:2
package ganrac

import __yyfmt__ "fmt"

//line parser.y:2

//line parser.y:5
type yySymType struct {
	yys  int
	node pNode
	num  int
}

const call = 57346
const list = 57347
const initvar = 57348
const name = 57349
const ident = 57350
const number = 57351
const f_true = 57352
const f_false = 57353
const all = 57354
const ex = 57355
const and = 57356
const or = 57357
const not = 57358
const abs = 57359
const plus = 57360
const minus = 57361
const comma = 57362
const mult = 57363
const div = 57364
const pow = 57365
const ltop = 57366
const gtop = 57367
const leop = 57368
const geop = 57369
const neop = 57370
const eqop = 57371
const assign = 57372
const eol = 57373
const lb = 57374
const rb = 57375
const lp = 57376
const rp = 57377
const lc = 57378
const rc = 57379
const unaryminus = 57380
const unaryplus = 57381

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"call",
	"list",
	"initvar",
	"name",
	"ident",
	"number",
	"f_true",
	"f_false",
	"all",
	"ex",
	"and",
	"or",
	"not",
	"abs",
	"plus",
	"minus",
	"comma",
	"mult",
	"div",
	"pow",
	"ltop",
	"gtop",
	"leop",
	"geop",
	"neop",
	"eqop",
	"assign",
	"eol",
	"lb",
	"rb",
	"lp",
	"rp",
	"lc",
	"rc",
	"unaryminus",
	"unaryplus",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:82
/*  start  of  programs  */

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 170

var yyAct = [...]int{
	2, 57, 32, 29, 35, 28, 19, 20, 21, 30,
	21, 33, 34, 62, 56, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	53, 15, 16, 57, 60, 17, 18, 55, 19, 20,
	21, 22, 23, 24, 25, 27, 26, 1, 61, 59,
	37, 12, 54, 52, 0, 15, 16, 0, 0, 17,
	18, 63, 19, 20, 21, 22, 23, 24, 25, 27,
	26, 0, 58, 15, 16, 36, 0, 17, 18, 0,
	19, 20, 21, 22, 23, 24, 25, 27, 26, 0,
	14, 9, 31, 7, 4, 5, 6, 9, 3, 7,
	4, 5, 6, 11, 10, 0, 0, 0, 0, 11,
	10, 17, 18, 0, 19, 20, 21, 13, 0, 8,
	0, 0, 0, 13, 0, 8, 15, 16, 0, 0,
	17, 18, 0, 19, 20, 21, 22, 23, 24, 25,
	27, 26, 15, 0, 0, 0, 17, 18, 0, 19,
	20, 21, 22, 23, 24, 25, 27, 26, 17, 18,
	0, 19, 20, 21, 22, 23, 24, 25, 27, 26,
}

var yyPact = [...]int{
	91, -1000, 59, -25, -1000, -1000, -1000, -31, 85, -32,
	85, 85, -1000, 42, -1000, 85, 85, 85, 85, 85,
	85, 85, 85, 85, 85, 85, 85, 85, 85, 85,
	17, -1000, 7, -13, -13, -19, -1000, -1000, 140, 128,
	-15, -15, -13, -13, -13, 93, 93, 93, 93, 93,
	93, 41, 14, 112, -1000, 13, -1000, 5, -1000, -1000,
	85, -1000, -1000, 112,
}

var yyPgo = [...]int{
	0, 53, 51, 4, 0, 47,
}

var yyR1 = [...]int{
	0, 5, 5, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 2, 2, 1,
	1, 3, 3,
}

var yyR2 = [...]int{
	0, 2, 4, 1, 1, 1, 1, 1, 3, 3,
	3, 4, 4, 3, 3, 3, 3, 3, 2, 2,
	3, 3, 3, 3, 3, 3, 1, 3, 2, 1,
	3, 1, 3,
}

var yyChk = [...]int{
	-1000, -5, -4, 7, 9, 10, 11, 8, 34, 6,
	19, 18, -2, 32, 31, 14, 15, 18, 19, 21,
	22, 23, 24, 25, 26, 27, 29, 28, 30, 34,
	-4, 7, 34, -4, -4, -3, 33, 8, -4, -4,
	-4, -4, -4, -4, -4, -4, -4, -4, -4, -4,
	-4, -4, -1, -4, 35, -3, 33, 20, 31, 35,
	20, 35, 8, -4,
}

var yyDef = [...]int{
	0, -2, 0, 7, 3, 4, 5, 6, 0, 0,
	0, 0, 26, 0, 1, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 7, 0, 18, 19, 0, 28, 31, 8, 9,
	13, 14, 15, 16, 17, 20, 21, 22, 23, 24,
	25, 0, 0, 29, 10, 0, 27, 0, 2, 11,
	0, 12, 32, 30,
}

var yyTok1 = [...]int{
	1,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:35
		{
			{
				yyytrace("gege")
			}
		}
	case 2:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.y:36
		{
			yyytrace("assign")
			stack.Push(yyDollar[2].node)
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:40
		{
			yyytrace("poly.num: " + yyDollar[1].node.str)
			stack.Push(yyDollar[1].node)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:41
		{
			yyytrace("true")
			stack.Push(yyDollar[1].node)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:42
		{
			yyytrace("false")
			stack.Push(yyDollar[1].node)
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:43
		{
			yyytrace("ident: " + yyDollar[1].node.str)
			stack.Push(yyDollar[1].node)
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:44
		{
			yyytrace("name: " + yyDollar[1].node.str)
			stack.Push(yyDollar[1].node)
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:45
		{
			yyytrace("and")
			stack.Push(yyDollar[2].node)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:46
		{
			yyytrace("or")
			stack.Push(yyDollar[2].node)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:47
		{
			yyVAL.node = yyDollar[2].node
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.y:48
		{
			yyytrace("call")
			stack.Push(newPNode(yyDollar[1].node.str, call, yyDollar[3].num, yyDollar[1].node.pos))
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.y:49
		{
			yyytrace("init")
			stack.Push(newPNode(yyDollar[1].node.str, initvar, yyDollar[3].num, yyDollar[1].node.pos))
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:50
		{
			yyytrace("+")
			stack.Push(yyDollar[2].node)
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:51
		{
			yyytrace("-")
			stack.Push(yyDollar[2].node)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:52
		{
			yyytrace("*")
			stack.Push(yyDollar[2].node)
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:53
		{
			yyytrace("/")
			stack.Push(yyDollar[2].node)
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:54
		{
			yyytrace("^")
			stack.Push(yyDollar[2].node)
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:55
		{
			yyytrace("-")
			stack.Push(newPNode("-.", unaryminus, 0, yyDollar[1].node.pos))
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:56
		{
			yyytrace("+.")
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:57
		{
			yyytrace("<")
			stack.Push(yyDollar[2].node)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:58
		{
			yyytrace(">")
			stack.Push(yyDollar[2].node)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:59
		{
			yyytrace("<=")
			stack.Push(yyDollar[2].node)
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:60
		{
			yyytrace(">=")
			stack.Push(yyDollar[2].node)
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:61
		{
			yyytrace("==")
			stack.Push(yyDollar[2].node)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:62
		{
			yyytrace("!=")
			stack.Push(yyDollar[2].node)
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:63
		{
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:67
		{
			yyytrace("list" + string(yyDollar[2].num))
			stack.Push(newPNode("_list", list, yyDollar[2].num, yyDollar[1].node.pos))
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:68
		{
			yyytrace("list0")
			stack.Push(newPNode("_list", list, 0, yyDollar[1].node.pos))
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:72
		{
			yyVAL.num = 1
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:73
		{
			yyVAL.num = yyDollar[1].num + 1
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:77
		{
			yyVAL.num = 1
			stack.Push(newPNode(yyDollar[1].node.str, ident, 0, yyDollar[1].node.pos))
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:78
		{
			yyVAL.num = yyDollar[1].num + 1
			stack.Push(newPNode(yyDollar[3].node.str, ident, 0, yyDollar[3].node.pos))
		}
	}
	goto yystack /* stack new state and value */
}
