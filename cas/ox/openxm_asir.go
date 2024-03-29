package ganrac

import (
	"fmt"
	. "github.com/hiwane/ganrac"
)

// j次 subresultant の DEG 次の係数を返す.
// DEG < 0 の場合は, subresultant を返す
// F = an x^n + ... + a1 x + a0
// =>
// Fs = [a0, a1, a2, ..., an, F]
const asir_init_str_sresj = `
def sresj(Fs, Gs, X, J, DEG) {
	local M, N, L, S, D, AI, BI, I;
    M = length(Fs) - 2;
    N = length(Gs) - 2;
	if (type(J) == 10) {
		J = int32ton(J);
	}
	if (type(DEG) == 10) {
		DEG = int32ton(DEG);
	}
	L = M+N-2*J;
	S = newmat(L, L);

	for (D = M; D >= 0; D--) {
		AI = Fs[D];
		for (I = 0; I < N - J && M-D+I < L-1; I++) {
			S[I][M-D+I] = AI;
		}
	}
	if (DEG >= 0) {
		for (I = N-J-1; I >= 0 && I-(N-J-1)+DEG >= 0; I--) {
			S[I][L-1] = Fs[I-(N-J-1)+DEG];
		}
	} else {
		for (I = N-J-1; I >= 0; I--) {
			S[I][L-1] = X^(N-J-1-I) * Fs[M+1];
		}
	}

	for (D = N; D >= 0; D--) {
		BI = Gs[D];
		for (I = 0; I < M - J && N-D+I < L-1; I++) {
			S[I+N-J][N-D+I] = BI;
		}
	}
	if (DEG >= 0) {
		for (I= M-J-1; I >= 0 && I-(M-J-1)+DEG >= 0; I--) {
			S[I+N-J][L-1] = Gs[I-(M-J-1)+DEG];
		}
	} else {
		for (I= M-J-1; I >= 0; I--) {
			S[I+N-J][L-1] = X^(M-J-1-I)* Gs[N+1];
		}
	}
	return det(S);
}
`

const asir_init_str_sres = `
def sres(F, G, X, CC) {
	FS = newvect(deg(F, X)+2);
	GS = newvect(deg(G, X)+2);
	for (I = 0; I < length(FS) - 1; I++) {
		FS[I] = coef(F, I, X);
	}
	FS[length(FS)-1] = F;
	for (I = 0; I < length(GS) - 1; I++) {
		GS[I] = coef(G, I, X);
	}
	GS[length(GS)-1] = G;

	K = length(FS);
	if (length(GS) < K) {
		K = length(GS);
	}
	K -= 2;
	if (type(CC) == 10) {
		CC = int32ton(CC);
	}

	RET = [];
	if (CC == 0) {
		for (J = K - 1; J >= 0; J--) {
			RET = cons(sresj(FS, GS, X, J, -1), RET);
		}
		return RET;
	}
	if (CC == 1) {
		for (J = K - 1; J >= 0; J--) {
			RET = cons(sresj(FS, GS, X, J, J), RET);
		}
		return RET;
	}
	if (CC == 2) {
		for (J = K - 1; J >= 0; J--) {
			RET = cons(sresj(FS, GS, X, J, 0), RET);
		}
		return RET;
	}
	if (CC == 3) {
		for (J = K - 1; J > 0; J--) {
			RET = cons(sresj(FS, GS, X, J, 0) + sresj(FS, GS, X, J, J) * X^J, RET);
		}
		J = 0;
		RET = cons(sresj(FS, GS, X, J, 0), RET);
		return RET;
	}
	return [];
}
`

// principal subresultant coefficient
const asir_init_str_psc = `
def psc(F, G, X, J) {
	FS = newvect(deg(F, X)+2);
	GS = newvect(deg(G, X)+2);
	for (I = 0; I < length(FS) - 1; I++) {
		FS[I] = coef(F, I, X);
	}
	FS[length(FS)-1] = F;
	for (I = 0; I < length(GS) - 1; I++) {
		GS[I] = coef(G, I, X);
	}
	GS[length(GS)-1] = G;
	return sresj(FS, GS, X, J, J);
}`
const asir_init_str_comb = `
def comb(A,B) {
	for (I=1, C=1; I<=B; I++) {
		C *= (A-I+1)/I;
	}
	return C;
}`
const asir_init_str_slope = `
def slope(F, G, X, K) {
    M = deg(F, X);
    N = deg(G, X);

	if (type(K) == 10) {
		K = int32ton(K);
	}
	L = N - K;
	S = newmat(L+1, L+1);
	CMK = comb(M, K+1);

	for (J = 0; J < L; J++) {
		S[0][J] = coef(F, M-J, X);
		for (I = 1; J+I < L; I++) {
			S[I][I+J] = S[0][J];
		}
		if (0 <= L-J && L-J <= L) {
			S[L-J][L] = (CMK - comb(M-J, K+1)) * S[0][J];
		}
		S[L][J] = coef(G, N-J, X);
	}
	J = L;
	S[0][L] = (CMK - comb(M-J, K+1)) * coef(F, M-J, X);
	S[L][L] = CMK * coef(G, K, X);
	return det(S);
}`

func (ox *OpenXM) Gcd(p, q *Poly) RObj {
	ox.ExecFunction("gcd", p, q)
	s, _ := ox.PopCMO()
	gob := ox.toGObj(s)
	return gob.(RObj)
}

func (ox *OpenXM) Factor(p *Poly) *List {
	// 因数分解
	err := ox.ExecFunction("fctr", p)
	if err != nil {
		ox.logger.Printf("Factor failed: input=%v, err=%s\n", p, err.Error())
		return nil
	}
	s, e := ox.PopCMO()
	if e != nil {
		ox.logger.Printf("Factor failed: err=%s\n", e.Error())
		ox.logger.Printf("Factor input: %v.\n", p)
		return nil
	}
	gob := ox.toGObj(s)
	return gob.(*List)
}

func (ox *OpenXM) Discrim(p *Poly, lv Level) RObj {
	dp := p.Diff(lv)
	ox.ExecFunction("res", NewPolyVar(lv), p, dp)
	qq, _ := ox.PopCMO()
	q := ox.toGObj(qq).(RObj)
	n := p.Deg(lv)
	if (n & 0x2) != 0 {
		q = q.Neg()
	}
	// 主係数で割る
	switch pc := p.Coef(lv, uint(n)).(type) {
	case *Poly:
		return q.(*Poly).Sdiv(pc)
	case NObj:
		return q.Div(pc)
	default:
		fmt.Printf("discrim: %v, pc=%v\n", p, pc)
	}
	return nil
}

func (ox *OpenXM) Resultant(p *Poly, q *Poly, lv Level) RObj {
	ox.ExecFunction("res", NewPolyVar(lv), p, q)
	qq, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("resultant %s\n", err.Error())
	} else if qq == nil {
		return NewInt(0)
	}
	return ox.toGObj(qq).(RObj)
}

func (ox *OpenXM) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	err := ox.ExecFunction("psc", p, q, NewPolyVar(lv), j)
	if err != nil {
		fmt.Printf("err: psc1 %s\n", err.Error())
	}
	qq, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("err: psc2 %s\n", err.Error())
		return nil
	} else if qq == nil {
		return NewInt(0)
	}
	return ox.toGObj(qq).(RObj)
}

func (ox *OpenXM) Sres(p *Poly, q *Poly, lv Level, cc int32) *List {
	err := ox.ExecFunction("sres", p, q, NewPolyVar(lv), cc)
	if err != nil {
		fmt.Printf("err: sres1 %s\n", err.Error())
	}
	qq, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("err: sres2 %s\n", err.Error())
		return nil
	} else if qq == nil {
		return NewList()
	}
	return ox.toGObj(qq).(*List)
}

func (ox *OpenXM) Slope(p *Poly, q *Poly, lv Level, k int32) RObj {
	err := ox.ExecFunction("slope", p, q, NewPolyVar(lv), k)
	if err != nil {
		fmt.Printf("err: slope1 %s\n", err.Error())
	}
	qq, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("err: slope2 %s\n", err.Error())
		return nil
	} else if qq == nil {
		return NewInt(0)
	}
	return ox.toGObj(qq).(RObj)
}

func (ox *OpenXM) GB(p *List, vars *List, n int) *List {
	// グレブナー基底
	var err error

	if n == 0 {
		err = ox.ExecFunction("nd_gr", p, vars, NewInt(0), NewInt(0))
	} else {
		// block order
		vn := NewList(
			NewList(NewInt(0), NewInt(int64(vars.Len()-n))),
			NewList(NewInt(0), NewInt(int64(n))))

		err = ox.ExecFunction("nd_gr", p, vars, NewInt(0), vn)
	}
	if err != nil {
		panic(fmt.Sprintf("gr failed: %v", err.Error()))
	}
	s, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("gr failed: %v", err.Error())
		return nil
	}

	gob := ox.toGObj(s)
	return gob.(*List)
}

func (ox *OpenXM) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {

	var err error

	if n == 0 {
		err = ox.ExecFunction("p_true_nf", p, gb, vars, NewInt(0))
	} else {
		// block order
		// https://www.asir.org/manuals/html-jp/man_172.html#SEC172
		vn := NewList(
			NewList(NewInt(0), NewInt(int64(vars.Len()-n))),
			NewList(NewInt(0), NewInt(int64(n))))

		err = ox.ExecFunction("p_true_nf", p, gb, vars, vn)
	}
	if err != nil {
		panic(fmt.Sprintf("p_nf failed: %v", err.Error()))
	}
	s, err := ox.PopCMO()
	if err != nil {
		fmt.Printf("p_nf failed: %v", err.Error())
		return nil, false
	}

	gob := ox.toGObj(s).(*List)

	m, _ := gob.Geti(1)
	mm, ok := m.(NObj)
	if !ok {
		panic(fmt.Sprintf("p_nf failed:\np=%v\ngb=%v\nret=%v\nden=%v", p, gb, gob, m))
	}
	sgn := mm.Sign() < 0

	m, _ = gob.Geti(0)

	return m.(RObj), sgn
}

func (ox *OpenXM) Eval(p string) (GObj, error) {
	ox.PushOxCMO(p)
	ox.PushOXCommand(SM_executeStringByLocalParser)
	s, err := ox.PopCMO()
	if err != nil {
		return nil, fmt.Errorf("popCMO failed %w", err)
	}
	gob := ox.toGObj(s)
	return gob, nil
}
