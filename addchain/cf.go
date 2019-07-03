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
//	[gencontfrac]            Amadou Tall and Ali Yassin Sanghare. Efficient computation of
//	                         addition-subtraction chains using generalized continued Fractions. Cryptology
//	                         ePrint Archive, Report 2013/466. 2013. https://eprint.iacr.org/2013/466
//	[hehcc:exp]              Christophe Doche. Exponentiation. Handbook of Elliptic and Hyperelliptic Curve
//	                         Cryptography, chapter 9. 2006.
//	                         https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch09.pdf

// ContinuedFractionStrategy is a method of choosing the auxiliary integer k in
// the continued fraction method outlined in [efficientcompaddchain].
type ContinuedFractionStrategy interface {
	// K returns values of k to try given n.
	K(n *big.Int) []*big.Int

	// Singleton returns whether every call to K will return one value of k. This
	// determines whether the resulting continued fractions sequence algorithm will
	// be logarithmic, and therefore suitable for large inputs.
	Singleton() bool

	fmt.Stringer
}

// ContinuedFractionStrategies lists all available continued fraction strategies.
var ContinuedFractionStrategies = []ContinuedFractionStrategy{
	BinaryStrategy{},
	CoBinaryStrategy{},
	DichotomicStrategy{},
	SqrtStrategy{},
	TotalStrategy{},
	DyadicStrategy{},
	FermatStrategy{},
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

// BinaryStrategy implements the binary strategy, which just sets k = floor(n/2). See [efficientcompaddchain] page 26.
// Since this is a singleton strategy it gives rise to a logarithmic sequence algoirithm that may not be optimal.
type BinaryStrategy struct{}

func (BinaryStrategy) String() string { return "binary" }

func (BinaryStrategy) Singleton() bool { return true }

// K returns floor(n/2).
func (BinaryStrategy) K(n *big.Int) []*big.Int {
	k := new(big.Int).Rsh(n, 1)
	return []*big.Int{k}
}

// CoBinaryStrategy implements the co-binary strategy, also referred to as the
// "modified-binary" strategy. See [efficientcompaddchain] page 26 or
// [gencontfrac] page 6. Since this is a singleton strategy it gives rise to a
// logarithmic sequence algoirithm that may not be optimal.
type CoBinaryStrategy struct{}

func (CoBinaryStrategy) String() string { return "co_binary" }

func (CoBinaryStrategy) Singleton() bool { return true }

// K returns floor(n/2) when n is even, or floor((n+1)/2) when n is odd.
func (CoBinaryStrategy) K(n *big.Int) []*big.Int {
	k := bigint.Clone(n)
	if k.Bit(0) == 1 {
		k.Add(k, bigint.One())
	}
	k.Rsh(k, 1)
	return []*big.Int{k}
}

// TotalStrategy returns all possible values of k less than n. This will result
// in the optimal continued fraction chain at a complexity of O(n^2 log^2(n)).
// Note that the optimal continued fraction chain is not necessarily the optimal
// chain. Must not be used for large inputs.
type TotalStrategy struct{}

func (TotalStrategy) String() string { return "total" }

// Singleton returns false, since the total strategy returns more than once k.
func (TotalStrategy) Singleton() bool { return false }

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

// DyadicStrategy implements the Dyadic Strategy, defined in
// [efficientcompaddchain] page 28. This gives rise to a sequence algorithm with
// complexity O(n*log^3(n)). Must not be used for large inputs.
type DyadicStrategy struct{}

func (DyadicStrategy) String() string { return "dyadic" }

// Singleton returns false, since the dyadic strategy returns more than once k.
func (DyadicStrategy) Singleton() bool { return false }

// K returns floor( n / 2ʲ ) for all j.
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

// FermatStrategy implements Fermat's Strategy, defined in
// [efficientcompaddchain] page 28. This returns a set of possible k of size
// O(log(log(n))), giving rise to a faster algorithm than the Dyadic strategy.
// This has been shown to be near optimal for small inputs. Must not be used for
// large inputs.
type FermatStrategy struct{}

func (FermatStrategy) String() string { return "fermat" }

// Singleton returns false, since Fermat's strategy returns more than once k.
func (FermatStrategy) Singleton() bool { return false }

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

// DichotomicStrategy is a singleton strategy, defined in
// [efficientcompaddchain] page 28. This gives rise to a logarithmic sequence
// algorithm, but the result is not necessarily optimal.
type DichotomicStrategy struct{}

func (DichotomicStrategy) String() string { return "dichotomic" }

// Singleton returns true, since the dichotomic strategy suggests just one k.
func (DichotomicStrategy) Singleton() bool { return true }

// K returns only one suggestion for k, namely floor( n / 2ʰ ) where h = log2(n)/2.
func (DichotomicStrategy) K(n *big.Int) []*big.Int {
	l := n.BitLen()
	h := uint(l) / 2
	k := new(big.Int).Div(n, bigint.Pow2(h))
	return []*big.Int{k}
}

// SqrtStrategy chooses k to be floor(sqrt(n)). See [gencontfrac] page 6. Since
// this is a singleton strategy, it gives rise to a logarithmic sequence
// algorithm that's not necessarily optimal.
type SqrtStrategy struct{}

func (SqrtStrategy) String() string { return "sqrt" }

// Singleton returns true, since the square root strategy suggests just one k.
func (SqrtStrategy) Singleton() bool { return false }

// K returns floor(sqrt(n)).
func (SqrtStrategy) K(n *big.Int) []*big.Int {
	sqrt := new(big.Int).Sqrt(n)
	return []*big.Int{sqrt}
}
