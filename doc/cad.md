# Cylindrical Algebraic Decomposition

## Projection Operator

| algorithm | implementation | citation |
| :-- | :--: | :--: |
| Collins' projection  | |
| [Hong's projection](../projhh.go) | ✔ |
| [McCallum's projection](../projmc.go) | ✔ |
| MC projection with equational constraints | |
| Lazard's projection |  |


```
> vars(a,b,s,t,u,x);
0
> F = ex([x], a*x+b != 0 && s*x^2+t*x+u <= 0);
ex([x], a*x+b!=0 && s*x^2+t*x+u<=0)
> cad(F); # McCallum Projection
go projalgo=0, lv=0
error: NOT well-oriented
> cad(F, 1); # Hong Projection
go projalgo=1, lv=0
(4*s*u-t^2<=0 && a*t-2*b*s!=0) || (a!=0 && s==0 && u==0) || (b!=0 && 4*s*u-t^2<0) || (b!=0 && u<=0) || (a!=0 && s<0) || a^2*u-a*b*t+b^2*s<0
```

```
> F = example("quad")[0]:
ex([w], x*w^2+y*w+z==0)
> C = cadinit(F);
CAD[x*w^2+y*w+z==0]
> cadproj(C);
[[x], [y], [z, 4*x*z-y^2], [x*w^2+y*w+z]]
> print(C, "proj");
[3, 0,i, 2] x*w^2+y*w+z
[2, 0, , 1] z
[2, 1, , 1] 4*x*z-y^2
[1, 0, , 1] y
[0, 0, , 1] x
 |  | |  |
 |  | |  |
 |  | |  +---  degree with respect to a main variable
 |  | +------  'i' if the polynomial is in the input formula
 |  +--------  index
 +-----------  level


> print(C, "proj", 2);
[2, 0, , 1] z
[2, 1, , 1] 4*x*z-y^2
> print(C, "proj", 2, 0);     # print(C, "proj", level, index) = print P(level, index)
[2, 0, , 1] z
coef[1]=+                 ...... coef(P(2, 0), 1) = c  where c > 0
discrim=+                 ...... discrim(P(2, 0)) = c  where c > 0
res[ 1]=+ P(1,  0)^2      ...... resultant(P(2,0), P(2,1)) = c * P(1, 0)^2 where c > 0
> print(C, "proj", 3, 0);
[3, 0,i, 2] x*w^2+y*w+z
coef[2]=+ P(0,  0)^1
coef[1]=+ P(1,  0)^1
coef[0]=+ P(2,  0)^1
discrim=- P(2,  1)^1
```

## Lifting

| algorithm | implementation | citation |
| :-- | :--: | :--: |
| [Symbolic-numeric CAD](../lift.go) | ✔| [1](https://www.sciencedirect.com/science/article/pii/S0304397512009413) |
| [Dynamic evaluation](../cad_de.go) | ✔| [1](https://dl.acm.org/doi/10.1006/jsco.1994.1057), [2](https://www.semanticscholar.org/paper/About-a-New-Method-for-Computing-in-Algebraic-Dora-Dicrescenzo/2ebef9590ca6ce106a45f491b0b864aa5a2206c2), [3](https://www.sciencedirect.com/science/article/pii/S0304397512009413) |
| Local projection | | [1](https://dl.acm.org/doi/10.1145/2608628.2608633) |

```
> F = example("quad")[0]:
> C = cadinit(F):
> cadproj(C);
[[x], [y], [z, 4*x*z-y^2], [x*w^2+y*w+z]]
> cadlift(C);
CAD[x*w^2+y*w+z==0]
> print(C, "sig");
sig(["sig"]) :: index=[], truth=-1
         ( )
  0,?, 3 (-) [-1.000000e+00,-1.000000e+00]
  1,?, 3 (1) [ 0.000000e+00, 0.000000e+00] 0
  2,?, 3 (+) [ 1.000000e+00, 1.000000e+00]
> print(C, "sig", 1);
sig(["sig" 1]) :: index=[1], truth=-1
         ( )
  0,?, 3 (-) [-1.000000e+00,-1.000000e+00]
  1,?, 3 (1) [ 0.000000e+00, 0.000000e+00] 0
  2,?, 3 (+) [ 1.000000e+00, 1.000000e+00]
> print(C, "sig", 1, 1);
sig(["sig" 1 1]) :: index=[1 1], truth=-1
         (   )
  0,f, 1 (- 0) [-1.000000e+00,-1.000000e+00]
  1,t, 1 (1 0) [ 0.000000e+00, 0.000000e+00] 0
  2,f, 1 (+ 0) [ 1.000000e+00, 1.000000e+00]
  | |  |  | |           |                    |
  | |  |  | |           |                    +--- defining polynomial
  | |  |  | |           +------------------------ isolating interval
  | |  |  | +------------------------------------ sign of P(2, 1) ... 0 if P(2, 1) vanish on cell(1, 1)          
  | |  |  +-------------------------------------- sign of P(2, 0) ... P(2, 0) is zero on cell(1, 1, 1) with multiplicity 1
  | |  +----------------------------------------- number of children
  | +-------------------------------------------- truth value
  +---------------------------------------------- index



> print(C, "cell", 1, 1);
--- information about the cell [1 1] 0xc0002b4fa0 ---
lv=1:y, de=false, exdeg=1, truth=-1 sgn=0
# of children=3
def.value    =0
signature    =(1)
iso.intv     =[0,0]
             =[0.000000e+00,0.000000e+00] = 0.000000e+00

> print(C, "cell", 1, 1, 1);
--- information about the cell [1 1 1] 0xc0002ce0a0 ---
lv=2:z, de=false, exdeg=1, truth=1 sgn=0
# of children=1
def.value    =0
signature    =(1 0)
iso.intv     =[0,0]
             =[0.000000e+00,0.000000e+00] = 0.000000e+00
```


## Soluation Formula Construction

- [Solution formula construction](https://dl.acm.org/doi/10.5555/929495)


## Demo

![cad](https://user-images.githubusercontent.com/7787544/199652778-84d4a90b-4906-4962-ac51-71b625bd9043.gif)

