package addchain

import (
	"fmt"
	"math/big"

	"github.com/mmcloughlin/ec3/internal/bigint"
	"github.com/mmcloughlin/ec3/internal/bigints"
)

// References:
//
//	[efficientcompaddchain]  Bergeron, F., Berstel, J. and Brlek, S. Efficient computation of addition
//	                         chains. Journal de theorie des nombres de Bordeaux. 1994.
//	                         http://www.numdam.org/item/JTNB_1994__6_1_21_0
//	[hehcc:exp]              Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                         Cryptography, chapter 9. 2006.
//	                         https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

// ContinuedFractionStrategy is a method of choosing the auxiliary integer k in
// the continued fraction method outlined in [efficientcompaddchain].
type ContinuedFractionStrategy interface {
	fmt.Stringer
	K(n *big.Int) []*big.Int
}

// ContinuedFractions uses the continued fractions method for finding an
// addition chain [efficientcompaddchain].
type ContinuedFractions struct {
	strategy ContinuedFractionStrategy
}

func NewContinuedFractions(s ContinuedFractionStrategy) ContinuedFractions {
	return ContinuedFractions{
		strategy: s,
	}
}

func (a ContinuedFractions) String() string {
	return fmt.Sprintf("continued_fractions(%s)", a.strategy)
}

func (a ContinuedFractions) FindSequence(targets []*big.Int) (Chain, error) {
	bigints.Sort(targets)
	return a.chain(targets), nil
}

func (a ContinuedFractions) minchain(n *big.Int) Chain {
	if bigint.IsPow2(n) {
		return bigint.Pow2UpTo(n)
	}

	if bigint.EqualInt64(n, 3) {
		return bigints.Int64s(1, 2, 3)
	}

	var min Chain
	for _, k := range a.strategy.K(n) {
		c := a.chain([]*big.Int{k, n})
		if min == nil || len(c) < len(min) {
			min = c
		}
	}

	return min
}

// chain produces a continued fraction chain for the given values. The slice ns
// must be in ascending order.
func (a ContinuedFractions) chain(ns []*big.Int) Chain {
	k := len(ns)
	if k == 1 || ns[k-2].Cmp(bigint.One()) <= 0 {
		return a.minchain(ns[k-1])
	}

	q, r := new(big.Int), new(big.Int)
	q.DivMod(ns[k-1], ns[k-2], r)

	cq := a.minchain(q)
	remaining := bigints.Clone(ns[:k-1])

	if bigint.IsZero(r) {
		return Product(a.chain(remaining), cq)
	}

	remaining = bigints.InsertSortedUnique(remaining, r)
	return Plus(Product(a.chain(remaining), cq), r)
}

type BinaryStrategy struct {
	Parity uint
}

func (b BinaryStrategy) String() string {
	if b.Parity == 0 {
		return "binary"
	}
	return "co_binary"
}

func (b BinaryStrategy) K(n *big.Int) []*big.Int {
	k := new(big.Int).Add(n, big.NewInt(int64(b.Parity)))
	k.Rsh(k, 1)
	return []*big.Int{k}
}

type TotalStrategy struct{}

func (TotalStrategy) String() string { return "total" }

// K returns {2,, 3, ..., n-1}.
func (TotalStrategy) K(n *big.Int) []*big.Int {
	ks := []*big.Int{}
	k := big.NewInt(2)
	one := bigint.One()
	for k.Cmp(n) < 0 {
		ks = append(ks, bigint.Clone(k))
		k.Add(k, one)
	}
	return ks
}

type DyadicStrategy struct{}

func (DyadicStrategy) String() string { return "dyadic" }

// K returns floor( n / 2^j ) for all j.
func (DyadicStrategy) K(n *big.Int) []*big.Int {
	ks := []*big.Int{}
	k := new(big.Int).Rsh(n, 1)
	one := bigint.One()
	for k.Cmp(one) > 0 {
		ks = append(ks, bigint.Clone(k))
		k.Rsh(k, 1)
	}
	return ks
}

type FermatStrategy struct{}

func (FermatStrategy) String() string { return "fermat" }

// K returns floor( n / 2^(2^j) ) for all j.
func (FermatStrategy) K(n *big.Int) []*big.Int {
	ks := []*big.Int{}
	k := new(big.Int).Rsh(n, 1)
	one := bigint.One()
	s := uint(1)
	for k.Cmp(one) > 0 {
		ks = append(ks, bigint.Clone(k))
		k.Rsh(k, s)
		s *= 2
	}
	return ks
}

type DichotomicStrategy struct{}

func (DichotomicStrategy) String() string { return "dichotomic" }

// K returns only one suggestion for k, namely floor( n / 2ʰ ) where h = log2(n)/2.
func (DichotomicStrategy) K(n *big.Int) []*big.Int {
	l := n.BitLen()
	h := uint(l) / 2
	k := new(big.Int).Div(n, bigint.Pow2(h))
	return []*big.Int{k}
}
